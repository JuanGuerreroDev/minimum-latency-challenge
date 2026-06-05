# Tech Stack Decisions — Minimum Latency System

---

## Decisiones Confirmadas

### 1. Lenguaje: Go 1.22+
- **Justificación**: Balance entre rendimiento, simplicidad, concurrencia nativa
- **Versión**: Go 1.22+ (última versión estable)
- **Módulos**: Go modules (`go.mod` + `go.sum`)

### 2. Framework de Networking: gnet v2
- **Package**: `github.com/panjf2000/gnet/v2`
- **Justificación**: Reactor Pattern nativo para Go, basado en event-loop, soporta IOCP en Windows
- **Alternativas consideradas**: evio (soporte Windows limitado), nbio (menos maduro), custom (demasiado esfuerzo)

### 3. Logging: log/slog (stdlib)
- **Justificación**: Incluido en Go stdlib desde 1.21, sin dependencias externas, structured logging nativo
- **Formato**: JSON output
- **Alternativas consideradas**: zerolog, zap — innecesarios dado que slog cubre los requerimientos

### 4. PBT Framework: rapid
- **Package**: `pgregory.net/rapid`
- **Justificación**: Lightweight, idiomatic Go, buena integración con `go test`, shrinking automático
- **Alternativas consideradas**: gopter (más features pero más complejo)
- **Compliance**: PBT-09

### 5. Comunicación: TCP sobre localhost
- **Protocolo**: Binario ultra-minimal (1 byte por mensaje)
- **Dirección**: `127.0.0.1:8080`
- **Conexión**: Persistente (sin reconnect)

---

## Go Runtime Configuration

### Configuración Actual: Default
El sistema se ejecuta con la configuración por defecto del Go runtime para establecer un baseline de rendimiento limpio.

```bash
# Configuración actual (default)
# No se establecen variables de entorno especiales
go run ./cmd/server
go run ./cmd/benchmark
```

### Opciones de Tuning para Futuras Optimizaciones

#### Opción A: Tuning Agresivo
```bash
# Desactivar GC completamente (evita GC pauses)
GOGC=off

# Single core (reduce context switching del Go scheduler)
GOMAXPROCS=1

# CPU affinity en Windows (fijar a un core específico)
# PowerShell:
$process = Start-Process -PassThru -FilePath ".\server.exe"
$process.ProcessorAffinity = 1  # Core 0

# Combinado:
$env:GOGC="off"; $env:GOMAXPROCS="1"; go run ./cmd/server
```

**Impacto esperado**: -30-50% latencia promedio, pero mayor uso de memoria (sin GC)
**Riesgo**: Memory leak si hay objetos que no se liberan; no recomendado para ejecución prolongada

#### Opción B: Tuning Moderado
```bash
# Desactivar GC (el benchmark es corto, no hay riesgo de memory leak)
GOGC=off

# Dejar GOMAXPROCS en default (usa todos los cores)
# Esto permite que gnet use su pool de event loops
$env:GOGC="off"; go run ./cmd/server
```

**Impacto esperado**: -20-30% latencia, sin overhead de GC durante el benchmark
**Riesgo**: Bajo para benchmarks cortos (10k iteraciones ≈ pocos segundos)

#### Opción C: Tuning Mínimo
```bash
# Solo reducir scheduling overhead
GOMAXPROCS=1

$env:GOMAXPROCS="1"; go run ./cmd/server
```

**Impacto esperado**: -10-15% latencia por reducción de context switching
**Riesgo**: Ninguno

#### Opción D: Memory Ballast (alternativa a GOGC=off)
```go
// En main(), pre-allocar un ballast grande para evitar GC triggers
ballast := make([]byte, 1<<30) // 1GB ballast
_ = ballast

// Esto hace que el GC piense que hay mucha memoria libre
// y ejecute menos frecuentemente, sin desactivarlo completamente
```

**Impacto esperado**: -15-25% GC frequency
**Riesgo**: Usa 1GB de memoria virtual (no física)

---

## Dependency Summary

| Package | Versión | Propósito | Tipo |
|---|---|---|---|
| `github.com/panjf2000/gnet/v2` | Latest stable | Reactor Pattern event loop | Runtime |
| `pgregory.net/rapid` | Latest stable | Property-based testing | Test only |
| Go stdlib (`net`, `time`, `os`, `sort`, `log/slog`) | Go 1.22+ | TCP, timing, I/O, logging | Runtime |

---

## Build & Run Commands

```bash
# Inicializar módulo Go
go mod init github.com/JuanGuerreroDev/minimum-latency-challenge

# Instalar dependencias
go get github.com/panjf2000/gnet/v2
go get pgregory.net/rapid  # solo para tests

# Build
go build -o server.exe ./cmd/server
go build -o benchmark.exe ./cmd/benchmark

# Run server
.\server.exe --port=8080

# Run benchmark (en otra terminal)
.\benchmark.exe --host=127.0.0.1 --port=8080 --iterations=10000

# Tests
go test ./...

# Benchmarks con allocations
go test -bench=. -benchmem ./internal/protocol/
go test -bench=. -benchmem ./internal/stats/

# PBT
go test -run TestProperty ./internal/protocol/
```
