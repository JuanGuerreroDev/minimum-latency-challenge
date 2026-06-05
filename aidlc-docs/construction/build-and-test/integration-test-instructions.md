# Integration Test Instructions — Minimum Latency System

## Purpose
Validar que el cliente de benchmark (`cmd/benchmark`) y el servidor reactor
(`cmd/server`) interoperan correctamente sobre TCP: el servidor responde `0x02`
a cada estímulo `0x01` dentro del mismo tick del event loop (BR-03).

## Test Scenarios

### Scenario 1: Benchmark Client → Server (round-trip happy path)
- **Description**: 10,000 estímulos sobre conexión TCP persistente; cada uno debe recibir respuesta.
- **Setup**: compilar `server.exe` y `benchmark.exe`.
- **Test Steps**:
  1. Arrancar el servidor: `.\server.exe --port=8080`
  2. En otra terminal: `.\benchmark.exe --host=127.0.0.1 --port=8080 --iterations=10000`
- **Expected Results**: `Successful iterations: 10000`, 0 errores, p99 < 1ms, `benchmark.log` creado con 10,000 trazas + resumen.
- **Cleanup**: detener el servidor (Ctrl+C → graceful shutdown, o `Stop-Process`).

### Scenario 2: Graceful Shutdown (SIGINT/SIGTERM)
- **Description**: el servidor cierra limpiamente al recibir señal de interrupción (Q2=A, PAT-RES-01).
- **Test Steps**: con el servidor corriendo, enviar Ctrl+C.
- **Expected Results**: logs `shutdown signal received` y `server shutdown complete`; el proceso termina sin panic.

### Scenario 3: Invalid Message Type (fail-safe)
- **Description**: un byte ≠ 0x01 debe ignorarse y loggearse, sin derribar el event loop (BR-01, PAT-RES-02).
- **Test Steps**: enviar un byte arbitrario (p.ej. `0x09`) con un cliente TCP manual.
- **Expected Results**: log `unknown message type`, sin respuesta, servidor sigue operativo.

## Setup Integration Test Environment

### 1. Start Required Services
```powershell
# Terminal 1
.\server.exe --port=8080
```

### 2. Configure Service Endpoints
El cliente apunta por defecto a `127.0.0.1:8080`; configurable con `--host`/`--port`.

## Run Integration Tests (script PowerShell usado)
```powershell
$srv = Start-Process .\server.exe "--port=8080" -PassThru `
  -RedirectStandardOutput server-out.log -RedirectStandardError server-err.log -WindowStyle Hidden
# polling de readiness sobre 127.0.0.1:8080 ...
.\benchmark.exe --host=127.0.0.1 --port=8080 --iterations=10000
Stop-Process -Id $srv.Id -Force
```

### Verify Service Interactions (resultado 2026-06-04)
- **Scenario 1**: ✅ 10,000/10,000 exitosas, 0 errores, p99 = 646.4µs.
- **Logs**: `server-out.log` (eventos slog JSON: starting/started/connection opened/closed), `benchmark.log` (trazas).

### 3. Cleanup
```powershell
Stop-Process -Name server -Force -ErrorAction SilentlyContinue
Remove-Item server-out.log, server-err.log -ErrorAction SilentlyContinue
```
