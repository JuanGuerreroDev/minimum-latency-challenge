# Minimum Latency Challenge

Sistema de comunicación TCP de **ultra-baja latencia**: ante un estímulo (1 byte
`0x01`), el servidor responde con `0x02` en el menor tiempo posible. Objetivo:
**latencia round-trip p99 < 1ms**.

Soporta dos modos de cliente: **benchmark** (10,000 iteraciones en bucle cerrado
para medición rigurosa) y **stimulus interactivo** (envío de estímulos ad-hoc en
cualquier momento presionando Enter).

Construido con la metodología **AI-DLC**. Diseño y artifacts en `aidlc-docs/`.

## Arquitectura

- **Reactor Pattern** sobre [`gnet/v2`](https://github.com/panjf2000/gnet) con
  **un único event loop** y `TCP_NODELAY` (Nagle deshabilitado).
- **Protocolo binario ultra-minimal**: 1 byte por mensaje, sin payload ni framing.
- **Zero-allocation hot path**: buffers pre-allocated en el servidor y los clientes.
- **Graceful shutdown** ante SIGINT/SIGTERM.
- Bind exclusivo a `127.0.0.1` (mínima superficie de ataque).

### Estructura del proyecto

```
minimum-latency-challenge/
├── cmd/
│   ├── server/main.go        # Entry point del servidor (gnet event loop)
│   ├── benchmark/main.go     # Cliente de benchmark (mide latencia round-trip)
│   └── stimulus/main.go      # Cliente interactivo (envío de estímulos ad-hoc)
├── internal/
│   ├── protocol/             # Encode/Decode binario 1 byte (+ PBT)
│   ├── stats/                # Cálculo de min/max/avg/p50/p95/p99 (+ PBT)
│   ├── logger/               # slog JSON + LatencyRecorder buffered
│   └── reactor/              # gnet EventHandler (hot path OnTraffic)
├── docs/
│   └── decisions/            # Architecture Decision Records (ADRs)
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
go build -o stimulus.exe ./cmd/stimulus

# 3. Arrancar el servidor (terminal 1)
.\server.exe --port=8080

# 4. Ejecutar el benchmark (terminal 2)
.\benchmark.exe --host=127.0.0.1 --port=8080 --iterations=10000

# 5. O enviar estímulos interactivos (terminal 2)
.\stimulus.exe --host=127.0.0.1 --port=8080 --log=stimulus.log
```

El benchmark imprime las estadísticas en stdout y escribe la trazabilidad por
petición en `benchmark.log`.

El stimulus interactivo permanece conectado al servidor y envía un estímulo cada
vez que se presiona Enter. Al salir (Ctrl+C o EOF), escribe el log de
trazabilidad con estadísticas resumen en el archivo configurado con `--log`.

## Tests

```powershell
# Tests unitarios y property-based
go test ./...

# Verificar zero-allocation en el hot path (debe reportar 0 allocs/op)
go test -bench=. -benchmem ./internal/protocol/
```

## Decisiones Arquitectónicas

Las decisiones de diseño del proyecto están documentadas como ADRs en
[`docs/decisions/`](docs/decisions/). Incluyen la elección del Reactor Pattern,
el protocolo binario, la estructura de binarios separados y el renombramiento
de `BenchmarkLogger` a `LatencyRecorder`.

## Opciones de tuning (post-baseline)

El servidor corre con configuración por defecto del runtime de Go para un
baseline limpio. Para optimizaciones adicionales (`GOGC=off`, `GOMAXPROCS=1`,
CPU affinity, memory ballast), ver
`aidlc-docs/construction/minimum-latency-system/nfr-requirements/tech-stack-decisions.md`.
