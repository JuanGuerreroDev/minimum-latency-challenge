# Documentación del Sistema — Minimum Latency Challenge

> Sistema de comunicación TCP de ultra-baja latencia. Ante un estímulo (1 byte
> `0x01`) el servidor responde `0x02` en el menor tiempo posible, con objetivo de
> latencia round-trip **p99 < 1 ms**. Diseñado y construido con la metodología
> **AI-DLC**.

---

## 1. Descripción Detallada de la Arquitectura

### 1.1 Visión general

El sistema consta de tres procesos que se comunican por TCP sobre localhost:

```
┌─────────────────────────┐         TCP (127.0.0.1:8080)        ┌──────────────────────────┐
│   cmd/benchmark          │   ── 0x01 (Stimulus) ──▶            │   cmd/server             │
│   (cliente de medición)  │                                     │   (gnet event loop)      │
│                          │   ◀── 0x02 (Response) ──            │                          │
│   mide round-trip ns     │                                     │   Reactor Handler        │
└─────────────────────────┘                                     └──────────────────────────┘

┌─────────────────────────┐         TCP (127.0.0.1:8080)        ┌──────────────────────────┐
│   cmd/stimulus           │   ── 0x01 (Stimulus) ──▶            │   cmd/server             │
│   (cliente interactivo)  │                                     │   (gnet event loop)      │
│                          │   ◀── 0x02 (Response) ──            │                          │
│   envío ad-hoc (Enter)   │                                     │   Reactor Handler        │
└─────────────────────────┘                                     └──────────────────────────┘
```

No hay infraestructura intermedia (colas, caches, balanceadores): cualquier hop
adicional incrementaría la latencia. La topología óptima para sub-milisegundo es
un único salto in-process por proceso. El servidor atiende múltiples conexiones
concurrentes — se puede correr benchmark y stimulus simultáneamente.

### 1.2 Patrón arquitectónico: Reactor Pattern

El servidor implementa el **Reactor Pattern** sobre la librería
[`gnet/v2`](https://github.com/panjf2000/gnet):

- **Synchronous Event Demultiplexer**: el event loop de gnet (`LC-01`) espera
  eventos de I/O del socket y los despacha.
- **Event Handler** (`internal/reactor`, `LC-02`): implementa los callbacks
  `OnBoot`, `OnOpen`, `OnTraffic`, `OnClose`, `OnShutdown`.
- **Single event loop** (`NumEventLoop=1`, `Multicore=false`): para una conexión
  1:1 elimina el overhead de coordinación entre loops y maximiza la localidad de
  caché del hot path.

### 1.3 Componentes

| Componente | Paquete | Rol | Path |
|---|---|---|---|
| Server entry | `cmd/server` | Arranca el event loop, graceful shutdown | hot/cold |
| Benchmark client | `cmd/benchmark` | Mide latencia round-trip (10k iteraciones) | hot |
| Stimulus client | `cmd/stimulus` | Envío interactivo de estímulos ad-hoc | hot |
| Reactor Handler | `internal/reactor` | `OnTraffic` responde al estímulo | hot |
| Protocol Codec | `internal/protocol` | Encode/Decode 1 byte | hot |
| Stats | `internal/stats` | min/max/avg/p50/p95/p99 | cold |
| Logger | `internal/logger` | slog JSON + LatencyRecorder buffered | cold |

Dependencias: `server → reactor, logger, gnet`; `benchmark → protocol, logger,
stats`; `stimulus → protocol, logger, stats`; `reactor → protocol, logger, gnet`;
los paquetes `protocol`, `stats` y `logger` son standalone.

### 1.4 Protocolo binario ultra-minimal

```
+--------+
| Type   |   Mensaje completo = 1 byte. Sin payload, sin framing.
| 1 byte |   0x01 = Stimulus (cliente→servidor)
+--------+   0x02 = Response (servidor→cliente)
```

Overhead de serialización: **0 bytes**. Se documentan alternativas extensibles
(payload "ping"/"pong", timestamp para latencia bidireccional) en
`aidlc-docs/.../business-rules.md`.

### 1.5 Decisiones de diseño orientadas a latencia

| Decisión | Justificación |
|---|---|
| **Single event loop** | Sin coordinación entre loops para 1 conexión |
| **TCP_NODELAY** (anti-Nagle) | Envío inmediato de cada byte; Nagle añadiría decenas de ms |
| **Zero-allocation hot path** | Buffer `[1]byte` pre-allocated; `conn.Next(-1)` sin copia; 0 allocs/op |
| **Bind localhost-only** | Mínima superficie de ataque; sin latencia de red externa |
| **Logging diferido** | Las trazas se acumulan en memoria; flush a disco post-benchmark (cero I/O en la medición) |
| **Runtime default** | Baseline limpio; tuning (GOGC/GOMAXPROCS) documentado para fases posteriores |

---

## 2. Cómo se Mide la Latencia

### 2.1 Definición

**Latencia round-trip** = tiempo transcurrido desde que el cliente envía el
estímulo (`0x01`) hasta que recibe la respuesta (`0x02`), medido sobre una única
conexión TCP persistente.

### 2.2 Método de medición

El cliente (`cmd/benchmark`) ejecuta un loop **síncrono single-shot** (BR-09):

```
para cada iteración (10,000):
    sendTime = time.Now()          // reloj monotónico de alta resolución
    conn.Write(stimulus[0x01])     // 1 byte
    conn.Read(readBuf)             // espera 1 byte de respuesta
    recvTime = time.Now()
    latency  = recvTime - sendTime
    registrar(sendTime, recvTime, latency)   // en memoria, sin I/O
```

Características que preservan la fidelidad de la medición:
- **`time.Now()`** de Go usa el reloj monotónico del SO (en Windows,
  QueryPerformanceCounter), inmune a ajustes del reloj de pared.
- **Sin I/O durante la medición**: las trazas se guardan en slices
  pre-allocated (`LatencyRecorder.Record`); el archivo `.log` se escribe solo al
  finalizar (BR-07).
- **Sin output intermedio** (BR-06): el reporte se imprime una sola vez al final.
- **Skip & continue** (BR-04): una iteración con error se cuenta y se omite; las
  estadísticas se calculan solo sobre iteraciones exitosas (BR-05).

### 2.3 Estadísticas calculadas

Sobre las latencias exitosas (ordenadas) se computan: `min`, `max`, `avg`,
`median (p50)`, `p95`, `p99`. Invariantes garantizadas y verificadas por
property-based testing: `Min ≤ Median ≤ Max` y `P50 ≤ P95 ≤ P99`.

### 2.4 Trazabilidad (log)

`benchmark.log` contiene una línea por petición con
`SendTime | RecvTime | Latency` (timestamps en RFC 3339 nano) más una sección de
resumen con las estadísticas. Es el registro de trazabilidad exigido por el
enunciado.

---

## 3. Justificación de Herramientas, Lenguajes y Metodologías

### 3.1 Lenguaje: Go 1.22+
- Concurrencia nativa, compilación a binario nativo, GC de baja pausa, y acceso a
  primitivas de red de bajo nivel sin la complejidad de C/C++.
- `time.Now()` monotónico de alta resolución y `log/slog` estructurado en stdlib.

### 3.2 Framework de red: gnet v2
- Implementación madura del **Reactor Pattern** para Go, basada en event loops
  con I/O no bloqueante; soporta IOCP en Windows.
- Evita el modelo goroutine-per-connection (más simple pero con más overhead de
  scheduling para latencia mínima).

### 3.3 Testing: rapid (property-based)
- Verifica invariantes sobre funciones puras (roundtrip de protocolo, invariantes
  de estadísticas) con shrinking automático y seeds reproducibles.

### 3.4 Logging: log/slog (stdlib)
- Structured logging JSON sin dependencias externas; usado solo en cold paths.

### 3.5 Metodología: AI-DLC
- Ciclo Inception → Construction → Operations con artifacts versionados en
  `aidlc-docs/` y gates de aprobación por etapa: Requirements, Application Design,
  Functional Design, NFR Requirements, NFR Design, Code Generation, Build & Test.
- Trazabilidad completa requisito → patrón → componente → código.

---

## 4. Build & Run (resumen)

```powershell
go mod tidy
go build -o server.exe ./cmd/server
go build -o benchmark.exe ./cmd/benchmark
go build -o stimulus.exe ./cmd/stimulus
.\server.exe --port=8080                                   # terminal 1
.\benchmark.exe --host=127.0.0.1 --port=8080 --iterations=10000   # terminal 2
.\stimulus.exe --host=127.0.0.1 --port=8080 --log=stimulus.log    # o modo interactivo
```

Resultados detallados en `docs/results-report.md`.
