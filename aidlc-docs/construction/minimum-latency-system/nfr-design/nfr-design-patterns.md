# NFR Design Patterns — Minimum Latency System

> Incorpora los NFR Requirements al diseño de la unidad mediante patrones de diseño concretos.
> **Decisiones de diseño base** (del plan): Q1=A → **Single event loop**, Q2=A → **Graceful shutdown habilitado**.

---

## 1. Performance Patterns

### PAT-PERF-01: Reactor Pattern (Single Event Loop)

**NFR cubierto**: NFR-PERF-01 (p99 < 1ms), RF-08 (Reactor Pattern)

**Decisión**: Un único event loop de gnet (`Options{NumEventLoop: 1}`).

**Justificación**:
- El benchmark usa **una conexión TCP persistente** (single-shot, 10k iteraciones secuenciales). Un solo loop elimina el overhead de scheduling entre múltiples loops y el cache-line bouncing entre cores.
- Mantiene el hot path en un único goroutine/thread del event loop, maximizando la localidad de caché (L1/L2 caliente para el handler y el buffer de respuesta).

**Configuración gnet**:
```go
gnet.Run(handler, "tcp://127.0.0.1:8080",
    gnet.WithMulticore(false),     // Single event loop (Q1=A)
    gnet.WithNumEventLoop(1),      // 1 loop explícito
    gnet.WithReusePort(false),     // No necesario con 1 loop
    gnet.WithTCPNoDelay(gnet.TCPNoDelay), // Disable Nagle → minimiza latencia
)
```

**Trade-off (Primera Ley de Arquitectura)**:
- ✅ Latencia mínima para 1 conexión, sin overhead de coordinación entre loops.
- ❌ No escala a múltiples conexiones concurrentes (aceptable: el escenario es un benchmark local 1:1).

---

### PAT-PERF-02: Disable Nagle's Algorithm (TCP_NODELAY)

**NFR cubierto**: NFR-PERF-01 (latencia round-trip)

**Decisión**: Activar `TCP_NODELAY` en el listener y en el cliente del benchmark.

**Justificación**: El algoritmo de Nagle agrupa paquetes pequeños esperando ACKs, añadiendo decenas de ms de latencia. Con mensajes de 1 byte esto sería catastrófico para p99. Deshabilitarlo fuerza el envío inmediato de cada byte.

```go
// Cliente benchmark
tcpConn := conn.(*net.TCPConn)
tcpConn.SetNoDelay(true)
```

---

### PAT-PERF-03: Zero-Allocation Hot Path

**NFR cubierto**: NFR-PERF-02 (0 allocs/op en hot path)

**Patrón**: **Pre-allocated buffers + reuse**. Ninguna asignación de heap en `OnTraffic`, `Encode` ni `Decode`.

**Técnicas aplicadas**:

| Técnica | Aplicación |
|---|---|
| Buffer de respuesta fijo en el struct | `response [1]byte` como campo del handler, reutilizado en cada `OnTraffic` |
| Lectura sin copia | `conn.Peek(1)` / `conn.Next(1)` de gnet en vez de `ReadAll` que asigna slices |
| Escritura directa | `conn.Write(h.response[:])` sobre el buffer pre-allocated |
| Slices pre-dimensionados en el cliente | `durations := make([]time.Duration, 0, iterations)` (capacidad reservada) |
| Sin interfaces ni boxing en el hot path | Tipos concretos (`byte`, `time.Duration`), evitar `interface{}` |

```go
type reactorHandler struct {
    gnet.BuiltinEventEngine
    response [1]byte // pre-allocated, reutilizado — 0 allocs
    logger   *slog.Logger
    eng      gnet.Engine
}

func (h *reactorHandler) OnTraffic(c gnet.Conn) gnet.Action {
    buf, _ := c.Next(-1)            // sin allocación: vista al buffer interno de gnet
    if len(buf) == 0 {
        return gnet.None
    }
    if buf[0] == protocol.TypeStimulus {
        h.response[0] = protocol.TypeResponse
        _ = c.AsyncWrite(h.response[:], nil) // o c.Write para sync
    } else {
        h.logger.Info("unknown message type", "type", buf[0]) // cold path
    }
    return gnet.None
}
```

**Verificación**: `go test -bench=. -benchmem` debe reportar `0 allocs/op` para `Encode`, `Decode` y el simulado de `OnTraffic`.

---

### PAT-PERF-04: Minimal Wire Protocol

**NFR cubierto**: NFR-PERF-03 (1 byte por mensaje)

**Patrón**: Protocolo binario sin framing ni header de longitud. El tipo (1 byte) **es** el mensaje completo. Cero overhead de serialización.

```
+--------+
| Type   |   ← mensaje completo = 1 byte
| 1 byte |
+--------+
0x01 = Stimulus | 0x02 = Response
```

---

### PAT-PERF-05: Baseline Runtime Configuration

**NFR cubierto**: NFR-PERF-04 (Go runtime default)

**Decisión**: Ejecutar con configuración por defecto del runtime de Go para establecer un baseline limpio. Las opciones de tuning (`GOGC=off`, `GOMAXPROCS=1`, memory ballast) quedan documentadas en `tech-stack-decisions.md` como optimizaciones futuras post-baseline.

---

## 2. Resilience Patterns

### PAT-RES-01: Graceful Shutdown (SIGINT/SIGTERM)

**NFR cubierto**: NFR-REL-01, Component 1 (Server). **Decisión Q2=A**.

**Patrón**: Captura de señales del OS + apagado limpio del event loop.

**Justificación**: Permite cerrar conexiones y hacer flush de logs sin perder datos ni dejar el socket en `TIME_WAIT` colgado. En Windows, gnet escucha `os.Interrupt`.

```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    handler := newReactorHandler(logger)
    go func() {
        <-ctx.Done()
        logger.Info("shutdown signal received, stopping engine")
        _ = handler.eng.Stop(context.Background()) // cierra el event loop limpiamente
    }()

    err := gnet.Run(handler, "tcp://127.0.0.1:8080",
        gnet.WithMulticore(false),
        gnet.WithNumEventLoop(1),
        gnet.WithTCPNoDelay(gnet.TCPNoDelay))
    if err != nil {
        logger.Error("server terminated", "err", err)
    }
}
```

`OnShutdown(eng)` se usa para loggear el cierre y confirmar la liberación de recursos.

---

### PAT-RES-02: Fail-Safe Input Handling (skip & continue)

**NFR cubierto**: NFR-REL-01, NFR-REL-02, NFR-SEC-01 (SECURITY-05), SECURITY-15

**Patrón**: **Defensive validation + fail closed**.

- **Servidor**: mensajes con tipo ≠ `0x01` se ignoran y se loggean; nunca crashean el event loop ni envían respuesta (fail closed, no fail open).
- **Cliente benchmark**: errores de `Write`/`Read` individuales incrementan un contador y ejecutan `continue` (skip & continue). Las estadísticas se calculan solo sobre iteraciones exitosas; el total de errores se reporta al final.

```go
// Benchmark loop — error resilience (NFR-REL-02)
for i := 0; i < iterations; i++ {
    sendTime := time.Now()
    if _, err := conn.Write(stimulus[:]); err != nil { errors++; continue }
    if _, err := conn.Read(readBuf[:]);  err != nil { errors++; continue }
    recvTime := time.Now()
    durations = append(durations, recvTime.Sub(sendTime))
}
```

---

### PAT-RES-03: Panic Recovery (cold path)

**NFR cubierto**: NFR-SEC-02 (SECURITY-15, global error handler)

**Patrón**: `recover()` en los límites del handler para que un panic inesperado no derribe el proceso. Se ubica **fuera del hot path** (solo se arma una vez, sin coste por iteración) para no afectar la latencia.

---

## 3. Security Patterns

### PAT-SEC-01: Minimal Attack Surface (localhost-only bind)

**NFR cubierto**: NFR-SEC-04 (SECURITY-09)

**Decisión**: Bind exclusivo a `127.0.0.1:8080`, nunca `0.0.0.0`. El servicio no es alcanzable desde la red externa.

### PAT-SEC-02: No Information Leakage

**NFR cubierto**: NFR-SEC-04 (SECURITY-09)

Las respuestas en el wire son únicamente `0x02`; nunca se serializan stack traces, errores internos ni metadata. Los errores van solo al log estructurado local.

### PAT-SEC-03: Dependency Pinning

**NFR cubierto**: NFR-SEC-03 (SECURITY-10)

`go.mod` con versión exacta de `gnet/v2` y `go.sum` committeado en control de versiones.

---

## 4. Observability Patterns

### PAT-OBS-01: Structured Logging (slog/JSON)

**NFR cubierto**: NFR-LOG-01 (SECURITY-03)

`log/slog` con handler JSON. Campos: `timestamp` (ISO 8601), `level`, `msg`. Sin datos sensibles. **El logging del servidor ocurre solo en cold paths** (`OnBoot`, `OnOpen`, `OnClose`, mensaje inválido), nunca en el éxito del hot path.

### PAT-OBS-02: Deferred Buffered Trace Logging

**NFR cubierto**: NFR-LOG-02

**Patrón**: **Record-in-memory, flush-on-exit**. Durante el benchmark, cada medición se escribe en slices pre-allocated (zero I/O en el loop de medición → zero impacto en latencia). Al finalizar, `FlushToFile` vuelca todo a `benchmark.log` con `bufio.Writer` y añade la sección de estadísticas.

```
Hot loop  →  Record() en slices pre-allocated (sin I/O)
   │
   ▼ (post-benchmark)
FlushToFile()  →  bufio.Writer  →  benchmark.log + resumen stats
```

---

## 5. Testing Patterns

### PAT-TEST-01: Property-Based Testing (partial)

**NFR cubierto**: NFR-TEST-01 (rapid, PBT-09)

Propiedades sobre funciones puras y roundtrips:
- **Encode/Decode roundtrip**: `Decode(Encode(t)) == t` para todo tipo válido.
- **Stats invariantes**: `Min ≤ Median ≤ Max` y `P50 ≤ P95 ≤ P99`.
- **Protocol size**: `Encode` siempre escribe exactamente 1 byte.

Shrinking habilitado; seed loggeado ante fallo para reproducibilidad.

### PAT-TEST-02: Benchmark Allocation Guards

**NFR cubierto**: NFR-TEST-03 (NFR-PERF-02)

`go test -bench -benchmem` como guardia de regresión: cualquier `allocs/op > 0` en el hot path falla la verificación de diseño.

---

## Resumen de Trazabilidad NFR → Patrón

| NFR | Patrón(es) |
|---|---|
| NFR-PERF-01 (p99 < 1ms) | PAT-PERF-01, PAT-PERF-02 |
| NFR-PERF-02 (0 allocs) | PAT-PERF-03, PAT-TEST-02 |
| NFR-PERF-03 (1 byte) | PAT-PERF-04 |
| NFR-PERF-04 (default runtime) | PAT-PERF-05 |
| NFR-REL-01 / NFR-REL-02 | PAT-RES-01, PAT-RES-02 |
| NFR-LOG-01 / NFR-LOG-02 | PAT-OBS-01, PAT-OBS-02 |
| NFR-SEC-01..04 | PAT-RES-02, PAT-RES-03, PAT-SEC-01..03 |
| NFR-TEST-01..03 | PAT-TEST-01, PAT-TEST-02 |
