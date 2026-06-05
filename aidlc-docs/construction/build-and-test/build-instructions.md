# Build Instructions — Minimum Latency System

## Prerequisites
- **Build Tool**: Go toolchain 1.22+ (verificado con **go1.26.4 windows/amd64**)
- **Dependencies** (resueltas vía `go mod tidy`):
  - `github.com/panjf2000/gnet/v2` v2.6.0 (runtime)
  - `pgregory.net/rapid` v1.1.0 (test only)
  - Indirectas: `golang.org/x/sys`, `golang.org/x/sync`, `go.uber.org/zap`, etc.
- **Environment Variables**: ninguna requerida (runtime default; ver tuning en tech-stack-decisions.md)
- **System Requirements**: Windows 10/11 x64; ~50 MB disco; cualquier CPU moderna

## Build Steps

### 1. Install Dependencies
```powershell
go mod tidy
```
Genera/actualiza `go.sum` con los hashes pinned (NFR-SEC-03).

### 2. Configure Environment
No se requiere configuración especial. El servidor escucha en `127.0.0.1:8080` por defecto.

### 3. Build All Units
```powershell
go build -o server.exe ./cmd/server
go build -o benchmark.exe ./cmd/benchmark
```

### 4. Verify Build Success
- **Expected Output**: sin errores, exit code 0.
- **Build Artifacts**:
  - `server.exe` (~6.7 MB)
  - `benchmark.exe` (~3.9 MB)
- **Common Warnings**: ninguno. `go vet ./...` debe quedar limpio.

## Troubleshooting

### Build Fails with Dependency Errors
- **Causa**: red sin acceso al proxy de módulos de Go, o `go.sum` corrupto.
- **Solución**: `go clean -modcache; go mod tidy`.

### Build Fails with Compilation Errors
- **Causa**: versión de Go < 1.21 (sin `log/slog`).
- **Solución**: instalar Go 1.22+ y verificar `go version`.

### `go` no se reconoce en la terminal
- **Causa**: PATH de la sesión no actualizado tras instalar Go.
- **Solución**: abrir una nueva terminal, o usar la ruta completa `& "C:\Program Files\Go\bin\go.exe"`.
