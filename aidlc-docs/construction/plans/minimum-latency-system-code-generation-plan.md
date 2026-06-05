# Code Generation Plan — minimum-latency-system

> **Single source of truth para la fase de Code Generation.** La generación (Parte 2) ejecuta EXACTAMENTE estos pasos en orden, sin lógica improvisada.

## Contexto de la Unidad

- **Unit**: `minimum-latency-system` (greenfield, single unit)
- **Project type**: Greenfield, Go 1.22+
- **Workspace root**: `C:\Users\ASUS\OneDrive\Documentos\Estudio\DiplomadoArquitecturaSoftware\Modulo1\DesarrolloActividad\minimum-latency-challenge`
- **Module path**: `github.com/JuanGuerreroDev/minimum-latency-challenge` (de tech-stack-decisions.md)
- **Code location**: workspace root (NUNCA `aidlc-docs/`)
- **Documentación de code-summaries**: `aidlc-docs/construction/minimum-latency-system/code/`

### Dependencias entre componentes (orden de generación bottom-up)
```
internal/protocol  (standalone)
internal/stats     (standalone)
internal/logger    (standalone) + internal/logger BenchmarkLogger (usa stats)
internal/reactor   → protocol, logger, gnet
cmd/server         → reactor, logger, gnet
cmd/benchmark      → protocol, logger, stats, stdlib
```

### Trazabilidad (Business Rules / NFR / Patterns)
| Artefacto de código | Reglas / NFR / Patrones |
|---|---|
| `internal/protocol/protocol.go` | BR-01, BR-02, PAT-PERF-03, PAT-PERF-04 |
| `internal/stats/stats.go` | BR-05, NFR-LOG-02, PAT-TEST-01 (invariantes) |
| `internal/logger/logger.go` | NFR-LOG-01, PAT-OBS-01 |
| `internal/logger/benchmark_logger.go` | BR-07, NFR-LOG-02, PAT-OBS-02, PAT-PERF-03 |
| `internal/reactor/handler.go` | BR-01, BR-03, PAT-PERF-01/02/03, PAT-RES-02/03, PAT-SEC-02 |
| `cmd/server/main.go` | PAT-PERF-01/02, PAT-RES-01 (graceful shutdown Q2=A), PAT-SEC-01 |
| `cmd/benchmark/main.go` | BR-04..BR-09, PAT-PERF-02/03, PAT-RES-02 |

---

## PLAN CHECKLIST

### Step 1 — Project Structure Setup (greenfield)
- [x] Crear `go.mod` con module path y `go 1.22`
- [x] Crear árbol de directorios: `cmd/server`, `cmd/benchmark`, `internal/{protocol,stats,logger,reactor}`

### Step 2 — Business Logic: `internal/protocol`
- [x] Generar `internal/protocol/protocol.go` (constantes `TypeStimulus`/`TypeResponse`/`MessageSize`, `Encode`, `Decode`, `ErrEmptyMessage`) — zero-allocation, payload vacío
- [x] Generar `internal/protocol/protocol_test.go` — PBT roundtrip + invariante de tamaño (rapid) + bench `-benchmem`

### Step 3 — Business Logic: `internal/stats`
- [x] Generar `internal/stats/stats.go` (`Stats` struct, `Calculate`, `String`, `Report`)
- [x] Generar `internal/stats/stats_test.go` — PBT invariantes (Min≤Median≤Max, P50≤P95≤P99, Count) + ejemplos

### Step 4 — Business Logic: `internal/logger`
- [x] Generar `internal/logger/logger.go` (wrapper `log/slog` JSON: `New`, `Info`, `Error`)
- [x] Generar `internal/logger/benchmark_logger.go` (`BenchmarkLogger` con slices pre-allocated: `NewBenchmarkLogger`, `Record`, `FlushToFile`)
- [x] Generar `internal/logger/benchmark_logger_test.go` — roundtrip record→flush

### Step 5 — Business Logic: `internal/reactor`
- [x] Generar `internal/reactor/handler.go` (`ReactorHandler` con buffer `[1]byte` pre-allocated; `OnBoot/OnOpen/OnClose/OnShutdown/OnTraffic`; `NewReactorHandler`) — hot path zero-allocation

### Step 6 — API/Entry Layer: `cmd/server`
- [x] Generar `cmd/server/main.go` (flags `--port`; single event loop `NumEventLoop=1`, `Multicore=false`, `TCPNoDelay`; bind `127.0.0.1`; graceful shutdown SIGINT/SIGTERM)

### Step 7 — API/Entry Layer: `cmd/benchmark`
- [x] Generar `cmd/benchmark/main.go` (flags `--host/--port/--iterations/--log`; `SetNoDelay(true)`; `runBenchmark` skip&continue; reporte stdout al final; flush a `benchmark.log`)

### Step 8 — Dependency Resolution
- [x] Ejecutar `go mod tidy` — resolvió `gnet/v2` v2.6.0 y `rapid` v1.1.0, generó `go.sum` (NFR-SEC-03). Go 1.26.4 instalado tras la generación inicial.

### Step 9 — Documentation Generation
- [x] Generar/actualizar `README.md` (build & run, estructura)
- [x] Generar code-summary en `aidlc-docs/construction/minimum-latency-system/code/code-summary.md`
- [x] **Nota**: `docs/system-documentation.md` y `docs/results-report.md` se completan en Build & Test (requieren resultados reales del benchmark)

### Step 10 — Validate Completeness
- [x] Verificar que todos los archivos del project structure existen (11 archivos de código/build + README confirmados)
- [x] Verificar imports coherentes con component-dependency.md
- [x] (La compilación y ejecución de tests ocurre en la fase **Build & Test**, no aquí)

---

## Archivos a Generar (resumen)

| # | Archivo | Tipo |
|---|---|---|
| 1 | `go.mod` | Build |
| 2 | `internal/protocol/protocol.go` | Código |
| 3 | `internal/protocol/protocol_test.go` | Test (PBT) |
| 4 | `internal/stats/stats.go` | Código |
| 5 | `internal/stats/stats_test.go` | Test (PBT) |
| 6 | `internal/logger/logger.go` | Código |
| 7 | `internal/logger/benchmark_logger.go` | Código |
| 8 | `internal/logger/benchmark_logger_test.go` | Test |
| 9 | `internal/reactor/handler.go` | Código |
| 10 | `cmd/server/main.go` | Código (entry) |
| 11 | `cmd/benchmark/main.go` | Código (entry) |
| 12 | `go.sum` | Build (vía go mod tidy) |
| 13 | `README.md` | Doc |
| 14 | `aidlc-docs/construction/minimum-latency-system/code/code-summary.md` | Doc |

**Total estimado**: 11 archivos de código/build + 2 de documentación + `go.sum`.

## Notas de Alcance
- No hay capa de Repository ni Frontend (sistema sin persistencia ni UI).
- No hay migraciones de base de datos.
- Los tests se **generan** aquí pero se **ejecutan** en Build & Test.
- `system-documentation.md` y `results-report.md` (requeridos por el enunciado) se generan en Build & Test con datos reales de latencia.
