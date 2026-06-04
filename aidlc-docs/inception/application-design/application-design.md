# Application Design — Minimum Latency Challenge

## Resumen Ejecutivo

Sistema de estímulo-respuesta de ultra-baja latencia implementado en Go con Reactor Pattern vía `gnet`. El sistema consta de un servidor TCP basado en event-loop y un cliente de benchmark que mide latencia round-trip sobre 10,000 iteraciones usando un protocolo binario ultra-minimal.

---

## Decisiones de Diseño Consolidadas

| Aspecto | Decisión | Justificación |
|---|---|---|
| Estructura proyecto | Package-based (`cmd/`, `internal/`) | Separación de concerns, Go idiomatic |
| Protocolo binario | Ultra-minimal (1 byte tipo + payload) | Máxima velocidad, mínimo overhead |
| Ciclo benchmark | Single-shot (10k requests, report, close) | Simplicidad, mediciones directas |
| Binarios | Separados (`server` + `benchmark`) | Independencia, deploy independiente |
| Framework | gnet v2 | Reactor Pattern nativo, IOCP Windows, producción probada |
| Logging | Buffered write (flush al final) | Zero impacto en mediciones de latencia |

---

## Arquitectura del Sistema

```
┌─────────────────────────────────┐     TCP      ┌──────────────────────────────────┐
│         Benchmark Client        │  (localhost)  │            Server                │
│                                 │              │                                  │
│  ┌───────────┐  ┌────────────┐  │    Binary    │  ┌────────────┐  ┌────────────┐  │
│  │ Protocol  │  │   Stats    │  │   Protocol   │  │  Reactor   │  │  Protocol  │  │
│  │ Encode/   │──┤ Calculate  │  │ ◄──────────► │  │  Handler   │──┤  Encode/   │  │
│  │ Decode    │  │ p50/p95/p99│  │              │  │ (OnTraffic)│  │  Decode    │  │
│  └───────────┘  └────────────┘  │              │  └────────────┘  └────────────┘  │
│                                 │              │                                  │
│  ┌───────────┐  ┌────────────┐  │              │  ┌────────────┐                  │
│  │ Benchmark │  │   Logger   │  │              │  │   Logger   │                  │
│  │ Logger    │  │ (struct)   │  │              │  │  (struct)  │                  │
│  │ (buffered)│  │            │  │              │  │            │                  │
│  └───────────┘  └────────────┘  │              │  └────────────┘                  │
│                                 │              │                                  │
│  cmd/benchmark/main.go          │              │  cmd/server/main.go              │
└─────────────────────────────────┘              └──────────────────────────────────┘
                                                          │
                                                    gnet v2 Event Loop
                                                    (IOCP on Windows)
```

---

## Componentes (6)

### Standalone (sin dependencias internas)
1. **`internal/protocol`** — Protocolo binario ultra-minimal (encode/decode)
2. **`internal/logger`** — Logger estructurado + BenchmarkLogger buffered
3. **`internal/stats`** — Calculador de estadísticas de latencia

### Con dependencias
4. **`internal/reactor`** — Event handler de gnet (depende de protocol, logger)

### Ejecutables
5. **`cmd/server`** — Server TCP con gnet Reactor (depende de reactor, logger)
6. **`cmd/benchmark`** — Cliente benchmark single-shot (depende de protocol, logger, stats)

---

## Protocolo Binario

```
+--------+-------------------+
| Type   | Payload           |
| 1 byte | Variable (0-255)  |
+--------+-------------------+

Type 0x01 = Stimulus ("ping")  →  5 bytes total
Type 0x02 = Response ("pong")  →  5 bytes total
```

**Flujo**:
1. Client: `Encode(0x01, "ping")` → 5 bytes → TCP Write
2. Server: `OnTraffic()` → `Decode()` → `Encode(0x02, "pong")` → TCP Write
3. Client: TCP Read → `Decode()` → measure latency

---

## Flujo del Benchmark

```
1. Connect (TCP dial to server)
2. Pre-allocate buffers & BenchmarkLogger(10000)
3. Loop 10,000x:
   a. sendTime = time.Now()
   b. Write stimulus to TCP
   c. Read response from TCP
   d. recvTime = time.Now()
   e. latency = recvTime - sendTime
   f. Record(sendTime, recvTime, latency)
4. Calculate stats (min, max, avg, p50, p95, p99)
5. Print report to stdout
6. Flush log to benchmark.log
7. Close connection
```

---

## Dependencia Externa

| Package | Propósito | Usado por |
|---|---|---|
| `github.com/panjf2000/gnet/v2` | Reactor Pattern event loop | `cmd/server`, `internal/reactor` |
| Go stdlib (`net`, `time`, `os`, `sort`) | TCP client, timing, I/O, sorting | `cmd/benchmark`, `internal/stats` |

---

## Estructura del Proyecto

```
minimum-latency-challenge/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── benchmark/
│       └── main.go
├── internal/
│   ├── reactor/
│   │   └── handler.go
│   ├── protocol/
│   │   ├── protocol.go
│   │   └── protocol_test.go
│   ├── logger/
│   │   ├── logger.go
│   │   └── benchmark_logger.go
│   └── stats/
│       ├── stats.go
│       └── stats_test.go
├── docs/
│   ├── system-documentation.md
│   └── results-report.md
├── go.mod
├── go.sum
└── README.md
```

---

## Artifacts de Diseño Detallados

- [components.md](file:///c:/Repositories/Github/Personal/minimum-latency-challenge/aidlc-docs/inception/application-design/components.md) — Definiciones de componentes y responsabilidades
- [component-methods.md](file:///c:/Repositories/Github/Personal/minimum-latency-challenge/aidlc-docs/inception/application-design/component-methods.md) — Firmas de métodos y tipos I/O
- [services.md](file:///c:/Repositories/Github/Personal/minimum-latency-challenge/aidlc-docs/inception/application-design/services.md) — Flujos de orquestación y servicios
- [component-dependency.md](file:///c:/Repositories/Github/Personal/minimum-latency-challenge/aidlc-docs/inception/application-design/component-dependency.md) — Matriz de dependencias y grafo
