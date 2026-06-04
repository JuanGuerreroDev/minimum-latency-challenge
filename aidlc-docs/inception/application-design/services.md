# Servicios y Orquestación — Minimum Latency Challenge

## Visión General

El sistema no tiene una capa de servicios tradicional (no hay microservicios, no hay API gateway). La orquestación es simple: un servidor y un cliente que se comunican directamente. Sin embargo, documentamos los flujos de orquestación internos de cada binario.

---

## Servicio 1: Server Lifecycle Service

**Responsabilidad**: Orquestar el ciclo de vida del servidor gnet.

### Flujo de Orquestación

```
main() 
  → Parsear flags (port)
  → Crear Logger
  → Crear ReactorHandler(logger)
  → Arrancar gnet.Run() con ReactorHandler en tcp://0.0.0.0:{port}
  → Bloquear hasta señal de shutdown (SIGINT/SIGTERM)
  → gnet cleanup automático
```

### Interacciones
- `main` → `logger.New()` → `reactor.NewReactorHandler()` → `gnet.Run()`
- `gnet` → `ReactorHandler.OnBoot()` (startup)
- `gnet` → `ReactorHandler.OnTraffic()` (por cada evento I/O)
- `gnet` → `ReactorHandler.OnClose()` (cuando el client desconecta)

---

## Servicio 2: Benchmark Execution Service

**Responsabilidad**: Orquestar la ejecución completa del benchmark single-shot.

### Flujo de Orquestación

```
main()
  → Parsear flags (host, port, iterations)
  → Crear Logger y BenchmarkLogger(iterations)
  → Conectar TCP al servidor (net.Dial)
  → Pre-allocar buffers de envío y recepción
  → Ejecutar runBenchmark():
      → Loop 10,000 veces:
          → Capturar sendTime (time.Now())
          → Enviar estímulo binario (protocol.Encode)
          → Leer respuesta binaria (conn.Read)
          → Capturar recvTime (time.Now())
          → Calcular latency = recvTime - sendTime
          → Registrar en BenchmarkLogger.Record()
  → Calcular stats.Calculate(durations)
  → Imprimir stats.Report() a stdout
  → BenchmarkLogger.FlushToFile("benchmark.log", stats)
  → Cerrar conexión
```

### Interacciones
- `main` → `logger.New()` → `logger.NewBenchmarkLogger()`
- `main` → `net.Dial()` → `runBenchmark()`
- `runBenchmark` → `protocol.Encode()` / `protocol.Decode()` (hot path)
- `main` → `stats.Calculate()` → `stats.Report()`
- `main` → `BenchmarkLogger.FlushToFile()`

---

## Flujo de Datos End-to-End

```
Benchmark Client                    Server (gnet)
     |                                   |
     |  1. Encode(Stimulus, "ping")      |
     |  2. TCP Write ──────────────────► |
     |                                   | 3. OnTraffic() triggered
     |                                   | 4. Decode(data) → Stimulus
     |                                   | 5. Encode(Response, "pong")
     |  ◄────────────────── TCP Write    | 6. conn.Write()
     |  7. TCP Read                      |
     |  8. Decode(data) → Response       |
     |  9. Calculate latency             |
     |  10. Record in buffer             |
     |                                   |
     |  [repeat 10,000 times]            |
     |                                   |
     |  11. Calculate stats              |
     |  12. Print report                 |
     |  13. Flush log file               |
     |  14. Close connection             |
```

---

## Notas de Diseño

1. **No hay service discovery**: Conexión directa por IP:port
2. **No hay retry logic**: Single-shot, si falla se reporta error
3. **No hay connection pooling**: Una conexión TCP persistente reutilizada
4. **El hot path es síncrono**: Enviar → recibir → medir → repetir
5. **I/O de logging es post-benchmark**: Zero impacto en mediciones de latencia
