# NFR Requirements — Minimum Latency System

---

## 1. Performance Requirements

### NFR-PERF-01: Latencia Round-Trip
- **Métrica**: Latencia end-to-end (envío estímulo → recepción respuesta)
- **Criterio de aceptación**: **p99 < 1ms** (1,000 microsegundos)
- **Objetivo ideal**: p50 < 100µs
- **Medición**: 10,000 iteraciones sobre conexión TCP persistente en localhost

### NFR-PERF-02: Zero-Allocation Hot Path
- **Métrica**: Allocaciones de heap en el hot path (OnTraffic + Encode/Decode)
- **Criterio de aceptación**: 0 allocaciones por iteración en el hot path
- **Verificación**: `go test -benchmem` debe reportar 0 allocs/op para hot path functions

### NFR-PERF-03: Protocolo Minimal
- **Métrica**: Bytes por mensaje en wire
- **Criterio de aceptación**: Exactamente 1 byte por mensaje (tipo sin payload)
- **Overhead**: 0 bytes de overhead de serialización

### NFR-PERF-04: Go Runtime Configuration
- **Configuración**: Default (sin tuning)
- **Justificación**: Establece baseline de rendimiento con configuración estándar
- **Documentación**: Opciones de tuning documentadas en tech-stack-decisions.md para futuras optimizaciones

---

## 2. Reliability Requirements

### NFR-REL-01: Estabilidad del Servidor
- **Criterio**: El servidor no debe crashear durante 10,000 iteraciones de benchmark
- **Error handling**: Mensajes inválidos se ignoran y logguean (no provocan crash)
- **Compliance**: SECURITY-15 (fail-safe defaults)

### NFR-REL-02: Benchmark Error Resilience
- **Criterio**: El benchmark continúa ejecutándose ante fallos individuales (skip and continue)
- **Reporte**: Número de errores reportado al final
- **Estadísticas**: Calculadas solo sobre iteraciones exitosas

---

## 3. Logging & Observability Requirements

### NFR-LOG-01: Structured Logging
- **Framework**: `log/slog` (Go stdlib, disponible desde Go 1.21)
- **Formato**: JSON structured output
- **Campos obligatorios**: timestamp (ISO 8601), level, message
- **Compliance**: SECURITY-03 (application-level logging)
- **Restricción**: Sin datos sensibles en logs

### NFR-LOG-02: Benchmark Trace Log
- **Formato**: Archivo .log con trazabilidad por petición
- **Contenido por línea**: SendTime, RecvTime, Latency
- **Resumen**: Estadísticas (min, max, avg, p50, p95, p99)
- **Estrategia de escritura**: Buffered write post-benchmark (zero impacto en mediciones)

---

## 4. Security Requirements

### NFR-SEC-01: Input Validation (SECURITY-05)
- **Validación**: Tipo de mensaje debe ser 0x01 o 0x02
- **Acción**: Mensajes inválidos ignorados, loggeados
- **Aplicable**: Solo en el server handler (OnTraffic)

### NFR-SEC-02: Error Handling (SECURITY-15)
- **Global error handler**: Recovery en el server para errores no manejados
- **Fail closed**: En caso de error, no enviar respuesta (no fail open)
- **Resource cleanup**: Conexiones cerradas en error paths

### NFR-SEC-03: Dependency Pinning (SECURITY-10)
- **Lock file**: `go.sum` committeado en version control
- **Versiones exactas**: gnet v2 pinned a versión específica en `go.mod`

### NFR-SEC-04: Hardening (SECURITY-09)
- **Error messages**: No exponer stack traces o detalles internos en respuestas
- **Minimal surface**: Solo escuchar en localhost (127.0.0.1), no 0.0.0.0

---

## 5. Testing Requirements

### NFR-TEST-01: Property-Based Testing (PBT Partial)
- **Framework**: `pgregory.net/rapid` (Go PBT framework — PBT-09)
- **Scope**: Funciones puras y roundtrips de serialización (PBT-02, PBT-03)
- **Propiedades a testar**:
  - Encode/Decode roundtrip
  - Stats invariantes (Min ≤ Median ≤ Max, P50 ≤ P95 ≤ P99)
  - Protocol size invariante (siempre 1 byte)
- **Shrinking**: Habilitado (PBT-08)
- **Reproducibilidad**: Seed logging en caso de fallo (PBT-08)

### NFR-TEST-02: Example-Based Testing (PBT-10)
- **Complementario**: Tests de ejemplo para business-critical paths
- **Cobertura**: Protocol encode/decode, Stats calculation, BenchmarkLogger flush

### NFR-TEST-03: Benchmark Tests
- **`go test -bench`**: Benchmarks de Go para medir allocaciones y throughput de funciones individuales
- **Hot path verification**: Confirmar 0 allocaciones en Encode/Decode/OnTraffic

---

## 6. Maintainability Requirements

### NFR-MAINT-01: Project Structure
- **Patrón**: Package-based (`cmd/`, `internal/`)
- **Módulos**: Go modules con `go.mod`
- **Documentación**: README.md, docs/system-documentation.md, docs/results-report.md

### NFR-MAINT-02: Code Quality
- **Linter**: `golangci-lint` con configuración por defecto
- **Formato**: `gofmt` standard formatting
- **Documentación**: Godoc comments en todas las funciones exportadas

---

## Acceptance Criteria Summary

| Métrica | Criterio | Prioridad |
|---|---|---|
| Latencia p99 | < 1ms | ⭐⭐⭐ CRÍTICO |
| Latencia p50 | < 100µs (ideal) | ⭐⭐ ALTO |
| Hot path allocs | 0 allocs/op | ⭐⭐⭐ CRÍTICO |
| Mensaje wire size | 1 byte | ⭐⭐ ALTO |
| Server stability | 0 crashes en 10k iterations | ⭐⭐ ALTO |
| PBT roundtrip pass | 100% | ⭐⭐ ALTO |
| Structured logging | SECURITY-03 compliant | ⭐ MEDIO |
