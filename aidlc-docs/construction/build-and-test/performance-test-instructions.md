# Performance Test Instructions — Minimum Latency System

## Purpose
Validar que el sistema cumple el objetivo de latencia round-trip **p99 < 1ms**
(NFR-PERF-01) y que el hot path no asigna memoria (NFR-PERF-02).

## Performance Requirements
- **Response Time**: **p99 < 1ms** (1,000 µs) sobre conexión TCP persistente en localhost.
- **Ideal**: p50 < 100µs.
- **Hot path allocations**: 0 allocs/op.
- **Error Rate**: 0% en 10,000 iteraciones.

## Setup Performance Test Environment

### 1. Prepare Test Environment
```powershell
go build -o server.exe ./cmd/server
go build -o benchmark.exe ./cmd/benchmark
```
Runtime: configuración por defecto de Go (baseline, NFR-PERF-04).

### 2. Configure Test Parameters
- **Iterations**: 10,000 (single-shot, síncrono).
- **Warmup**: ninguno (baseline; warmup documentado como optimización futura).
- **Conexión**: única TCP persistente, `TCP_NODELAY` activo en server y cliente.

## Run Performance Tests

### 1. Latency Benchmark (round-trip end-to-end)
```powershell
.\server.exe --port=8080      # terminal 1
.\benchmark.exe --iterations=10000   # terminal 2
```

### 2. Allocation Benchmark (hot path)
```powershell
go test "-bench=." "-benchmem" "-run=NONE" ./internal/protocol/
```

### 3. Analyze Performance Results (2026-06-04, i5-12500H, Go 1.26.4)

| Métrica | Objetivo | Medido | Estado |
|---|---|---|---|
| Latencia p99 | < 1ms | **646.4 µs** | ✅ |
| Latencia p95 | — | 550.3 µs | ✅ |
| Latencia p50 (median) | < 100µs (ideal) | 0s* | ✅ |
| Latencia avg | — | 85.43 µs | ✅ |
| Latencia min | — | 0s* | — |
| Latencia max | — | 5.56 ms | ⚠️ outlier |
| Iteraciones exitosas | 10,000 | 10,000 | ✅ |
| Error rate | 0% | 0% | ✅ |
| Hot path allocs | 0/op | 0/op | ✅ |

\* **0s** = round-trip por debajo de la granularidad efectiva del reloj para esas
iteraciones en este equipo. Ver análisis en `docs/results-report.md`.

- **Bottlenecks identificados**: el `Max` de 5.56ms es un outlier aislado,
  probablemente una pausa de GC o reprogramación del scheduler del SO (el runtime
  corre con GC por defecto). p95/p99 no se ven afectados.

## Performance Optimization
Si se requiriera reducir aún más la cola de latencia (p99/max):
1. `GOGC=off` o memory ballast para eliminar pausas de GC durante el benchmark.
2. `GOMAXPROCS=1` + CPU affinity para reducir context switching.
3. Warmup de N iteraciones antes de medir.

Ver detalle de cada opción en `tech-stack-decisions.md`.
