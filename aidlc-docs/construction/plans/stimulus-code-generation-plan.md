# Code Generation Plan — Stimulus Client (Interactive Mode)

## Unit Context
- **Unit**: `cmd/stimulus` — Cliente interactivo de estímulos
- **Stories/Requirements**: RF-09 (estímulo interactivo), RF-10 (trazabilidad)
- **Dependencies**: `internal/protocol` (encode/decode), `internal/logger` (Logger + BenchmarkLogger)
- **Interface**: Protocolo binario existente (0x01 → 0x02), TCP localhost

## Code Location
- **Application Code**: `cmd/stimulus/main.go` (workspace root)
- **Documentation**: `aidlc-docs/construction/stimulus/code/` (markdown summary)

---

## Generation Steps

### Step 1: Create `cmd/stimulus/main.go`
- [x] Implementar cliente interactivo con:
  - Flag `--host` (default: 127.0.0.1)
  - Flag `--port` (default: 8080)
  - Flag `--log` (default: stimulus.log)
  - Conexión TCP persistente con TCP_NODELAY
  - Loop de stdin: leer línea (Enter) → enviar 0x01 → leer 0x02 → imprimir latencia
  - Buffers pre-allocated `[1]byte` para send y recv (zero-allocation hot path)
  - Graceful shutdown con SIGINT/SIGTERM (flush log, cerrar conexión)
  - Reutilizar `logger.BenchmarkLogger` para acumular mediciones
  - Flush del log al terminar (con estadísticas si hay >= 2 mediciones)

### Step 2: Verify Build
- [x] Ejecutar `go build ./cmd/stimulus` para confirmar compilación limpia

### Step 3: Create Code Summary Documentation
- [x] Crear `aidlc-docs/construction/stimulus/code/code-summary.md` con descripción del binario generado

---

## Quality Gates (from execution plan)
- [x] Binario compila sin errores → Step 2
- [ ] Envía estímulo y recibe respuesta correcta → Build and Test
- [ ] Muestra latencia en consola → Step 1 (print en stdout)
- [ ] Escribe log de trazabilidad → Step 1 (BenchmarkLogger + FlushToFile)
- [ ] Graceful shutdown con SIGINT → Step 1 (signal.NotifyContext)
- [ ] Zero-allocation en hot path → Step 1 (buffers pre-allocated fuera del loop)
