# Code Summary — Stimulus Client (Interactive Mode)

## Archivos Generados/Modificados

| Archivo | Tipo | Descripción |
|---|---|---|
| `cmd/stimulus/main.go` | Nuevo | Cliente TCP interactivo para envío de estímulos ad-hoc |
| `internal/logger/latency_recorder.go` | Renombrado | Antes `benchmark_logger.go` — Renombrado `BenchmarkLogger` → `LatencyRecorder` (SRP) |
| `internal/logger/latency_recorder_test.go` | Renombrado | Antes `benchmark_logger_test.go` — Referencias actualizadas |
| `internal/logger/logger.go` | Modificado | Package comment actualizado |
| `cmd/benchmark/main.go` | Modificado | Referencias actualizadas a `LatencyRecorder` |

## Diseño del Binario

### Flujo Principal
1. Parsea flags (`--host`, `--port`, `--log`)
2. Establece conexión TCP persistente con TCP_NODELAY
3. Registra handler de señales (SIGINT/SIGTERM) para graceful shutdown
4. Pre-allocates buffers de envío y recepción (`[1]byte`)
5. Inicia goroutine de lectura de stdin (scanner)
6. Por cada Enter: send 0x01 → recv 0x02 → print latencia → record en LatencyRecorder
7. Al recibir señal o EOF: flush log a disco con estadísticas (si >= 2 mediciones)

### Decisiones de Implementación

| Decisión | Justificación |
|---|---|
| Goroutine para stdin | `bufio.Scanner.Scan()` es blocking; necesitamos que `select` escuche señales en paralelo |
| `durations` slice local | `LatencyRecorder` no expone duraciones; tracking local evita modificar su API |
| Capacidad inicial 64 | Uso interactivo típico: decenas de estímulos, no miles. Crece con append si se excede |
| Reutilización de `LatencyRecorder` | Mismo formato de log que el benchmark — consistencia del proyecto |
| Rename `BenchmarkLogger` → `LatencyRecorder` | SRP: el nombre refleja la responsabilidad real ("registrar latencias") sin acoplar al benchmark |
| Sin modificación de API existente | Solo rename de símbolos, sin cambios de comportamiento |

### Cumplimiento de Quality Gates

| Gate | Estado | Evidencia |
|---|---|---|
| Binario compila sin errores | ✅ | `go build ./cmd/stimulus` exit 0 |
| Zero-allocation hot path | ✅ | Buffers `[1]byte` pre-allocated fuera del loop |
| Graceful shutdown | ✅ | `signal.NotifyContext` + flush en exit path |
| Log de trazabilidad | ✅ | `LatencyRecorder.FlushToFile` con flag `--log` |
| Muestra latencia en consola | ✅ | `fmt.Printf` con duración por cada estímulo |
| Tests existentes sin regresión | ✅ | `go test ./...` all pass |
| SRP respetado | ✅ | `LatencyRecorder` es genérico, usado por benchmark y stimulus |
