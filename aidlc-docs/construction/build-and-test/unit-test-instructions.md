# Unit Test Execution — Minimum Latency System

## Run Unit Tests

### 1. Execute All Unit Tests
```powershell
go test ./...
```

### 2. Review Test Results
- **Expected**: todos los paquetes con tests PASS, 0 fallos.
- **Paquetes con tests**:
  - `internal/protocol` — PBT roundtrip + invariante de tamaño (rapid) + ejemplos
  - `internal/stats` — PBT invariantes (Min≤Median≤Max, P50≤P95≤P99, Count, Median==P50) + ejemplos
  - `internal/logger` — roundtrip Record→FlushToFile + crecimiento sobre capacidad
- **Paquetes sin tests** (entry points / handler, validados vía integración): `cmd/server`, `cmd/benchmark`, `internal/reactor`.

### 3. Resultado de la ejecución (2026-06-04, Go 1.26.4)
```
?   cmd/benchmark         [no test files]
?   cmd/server            [no test files]
ok  internal/logger       0.981s
ok  internal/protocol     1.609s
?   internal/reactor      [no test files]
ok  internal/stats        1.613s
```
**Estado: ✅ PASS** (3/3 paquetes con tests).

### 4. Verificación de Zero-Allocation (NFR-PERF-02)
```powershell
go test "-bench=." "-benchmem" "-run=NONE" ./internal/protocol/
```
Resultado:
```
BenchmarkEncode-16   1000000000   0.1747 ns/op   0 B/op   0 allocs/op
BenchmarkDecode-16   1000000000   0.1810 ns/op   0 B/op   0 allocs/op
```
**Estado: ✅ 0 allocs/op** en el hot path del protocolo.

> Nota PowerShell: citar los flags regex (`"-bench=."`, `"-run=NONE"`) evita que el
> shell interprete `.` y `$` y rompa el parseo de paquetes.

### 5. Fix Failing Tests
Si algún test falla, `rapid` imprime el seed para reproducir:
```powershell
go test "-run=TestProperty..." "-rapid.seed=<seed>" ./internal/...
```
