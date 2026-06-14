# Requerimientos — Intent V2: Stimulus Client Mode

## Intent Analysis Summary

- **User Request**: Agregar cliente interactivo que permita enviar estímulos al servidor en cualquier momento, independiente del benchmark
- **Request Type**: Enhancement (cumplimiento completo de RF-02)
- **Scope Estimate**: Single Component (nuevo binario `cmd/stimulus`)
- **Complexity Estimate**: Simple
- **Depth**: Minimal

---

## Gap Identificado

El requerimiento RF-02 ("Al recibir un estímulo, el sistema debe retornar una respuesta") solo se cumple actualmente dentro del loop cerrado de `cmd/benchmark`. No existe forma de enviar un estímulo ad-hoc al servidor en cualquier momento. Este intent cierra ese gap.

---

## Requerimientos Funcionales (Nuevos)

### RF-09: Cliente de Estímulo Interactivo
El sistema debe proveer un binario independiente (`cmd/stimulus`) que:
- Se conecte al servidor TCP en localhost (conexión persistente)
- Permanezca en modo interactivo esperando input del usuario
- Envíe un estímulo (0x01) al servidor cada vez que el usuario presione Enter
- Muestre en consola la latencia round-trip de cada estímulo enviado
- Cierre la conexión y termine al recibir señal de interrupción (Ctrl+C) o EOF

### RF-10: Trazabilidad del Cliente Interactivo
El cliente de estímulo debe:
- Mostrar la latencia de cada estímulo en stdout (formato legible)
- Escribir un archivo de log (.log) configurable mediante flag `--log`
- Registrar por cada estímulo: timestamp de envío, timestamp de recepción, latencia calculada
- Hacer flush del log al terminar la sesión (graceful shutdown)

---

## Requerimientos No Funcionales (Aplicables)

Los siguientes RNF del intent original aplican sin cambios:
- **RNF-01**: Latencia objetivo < 1ms (aplica a cada estímulo individual)
- **RNF-02**: Go como lenguaje
- **RNF-03**: Windows como plataforma
- **RNF-04**: TCP sobre localhost
- **RNF-05**: Zero-allocation en hot path del cliente (reutilizar buffers)
- **RNF-07**: Logging estructurado

---

## Restricciones de Diseño

1. **Reutilizar paquetes existentes**: `internal/protocol` para encode/decode, `internal/logger` para logging
2. **No modificar** el servidor ni el benchmark existente
3. **Patrón de proyecto existente**: Nuevo binario en `cmd/stimulus/main.go`
4. **Conexión persistente**: Una sola conexión TCP durante toda la sesión interactiva
5. **Graceful shutdown**: Capturar SIGINT/SIGTERM para flush del log y cierre limpio

---

## Extension Configuration (Sin Cambios)

| Extension | Status | Enforcement Mode |
|---|---|---|
| Security Baseline | Enabled | Pragmático (proyecto experimental) |
| Property-Based Testing | Enabled | Partial (solo funciones puras y roundtrips de serialización) |

---

## Entregables Esperados

1. `cmd/stimulus/main.go` — Binario interactivo de envío de estímulos
2. Tests unitarios aplicables (PBT para protocol roundtrip ya existe)
3. Actualización de documentación si es necesario
