# Minimum Latency Challenge

Sistema de comunicación TCP de **ultra-baja latencia**: ante un estímulo (1 byte
`0x01`), el servidor responde con `0x02` en el menor tiempo posible. Objetivo:
**latencia round-trip p99 < 1ms**.

Construido con la metodología **AI-DLC**. Diseño y artifacts en `aidlc-docs/`.

## Arquitectura

- **Reactor Pattern** sobre [`gnet/v2`](https://github.com/panjf2000/gnet) con
  **un único event loop** y `TCP_NODELAY` (Nagle deshabilitado).
- **Protocolo binario ultra-minimal**: 1 byte por mensaje, sin payload ni framing.
- **Zero-allocation hot path**: buffers pre-allocated en el servidor y el cliente.
- **Graceful shutdown** ante SIGINT/SIGTERM.
- Bind exclusivo a `127.0.0.1` (mínima superficie de ataque).

### Estructura del proyecto

```
minimum-latency-challenge/
├── cmd/
│   ├── server/main.go        # Entry point del servidor (gnet event loop)
│   └── benchmark/main.go     # Cliente de benchmark (mide latencia round-trip)
├── internal/
│   ├── protocol/             # Encode/Decode binario 1 byte (+ PBT)
│   ├── stats/                # Cálculo de min/max/avg/p50/p95/p99 (+ PBT)
│   ├── logger/               # slog JSON + BenchmarkLogger buffered
│   └── reactor/              # gnet EventHandler (hot path OnTraffic)
├── go.mod
└── aidlc-docs/               # Artifacts de diseño AI-DLC
```

## Requisitos

- **Go 1.22+** ([descargar](https://go.dev/dl/))

## Build & Run

```powershell
# 1. Resolver dependencias (genera go.sum)
go mod tidy

# 2. Compilar
go build -o server.exe ./cmd/server
go build -o benchmark.exe ./cmd/benchmark

# 3. Arrancar el servidor (terminal 1)
.\server.exe --port=8080

# 4. Ejecutar el benchmark (terminal 2)
.\benchmark.exe --host=127.0.0.1 --port=8080 --iterations=10000
```

El benchmark imprime las estadísticas en stdout y escribe la trazabilidad por
petición en `benchmark.log`.

## Tests

```powershell
# Tests unitarios y property-based
go test ./...

# Verificar zero-allocation en el hot path (debe reportar 0 allocs/op)
go test -bench=. -benchmem ./internal/protocol/
```

## Opciones de tuning (post-baseline)

El servidor corre con configuración por defecto del runtime de Go para un
baseline limpio. Para optimizaciones adicionales (`GOGC=off`, `GOMAXPROCS=1`,
CPU affinity, memory ballast), ver
`aidlc-docs/construction/minimum-latency-system/nfr-requirements/tech-stack-decisions.md`.
