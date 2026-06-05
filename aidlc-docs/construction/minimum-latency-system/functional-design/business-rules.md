# Business Rules — Minimum Latency System

## Reglas de Validación del Protocolo

### BR-01: Validación de Tipo de Mensaje
- **Regla**: Solo los tipos `0x01` (Stimulus) y `0x02` (Response) son válidos
- **Acción en violación**: Ignorar silenciosamente, loggear con nivel INFO
- **Justificación**: No desperdiciar ciclos del event loop en mensajes inválidos

### BR-02: Tamaño de Mensaje
- **Regla**: Un mensaje válido tiene exactamente **1 byte** (solo el tipo, sin payload)
- **Acción en violación**: Procesar solo el primer byte, ignorar bytes restantes
- **Justificación**: Protocolo ultra-minimal para máxima velocidad

### BR-03: Respuesta Obligatoria a Estímulo
- **Regla**: Todo mensaje de tipo `0x01` (Stimulus) DEBE generar una respuesta de tipo `0x02` (Response)
- **Latencia máxima**: La respuesta debe enviarse dentro del mismo tick del event loop
- **Justificación**: Core del sistema — no se permite procesamiento diferido

---

## Reglas del Benchmark

### BR-04: Skip on Error
- **Regla**: Si una iteración individual falla (write error, read error, timeout), se registra el error y se continúa con la siguiente iteración
- **No se hace retry**: Cada iteración es un intento único
- **Registro**: El número total de errores se reporta al final

### BR-05: Estadísticas Solo Sobre Éxitos
- **Regla**: Las estadísticas de latencia (min, max, avg, percentiles) se calculan SOLO sobre iteraciones exitosas
- **Justificación**: Las iteraciones fallidas distorsionarían las métricas de latencia

### BR-06: Reporte al Finalizar
- **Regla**: El reporte de resultados se imprime SOLO al finalizar todas las iteraciones
- **No hay output intermedio**: El benchmark es silencioso durante la ejecución para minimizar overhead

### BR-07: Logging Post-Benchmark
- **Regla**: El archivo .log se escribe DESPUÉS de completar todas las iteraciones
- **Justificación**: Cualquier I/O durante el benchmark afectaría las mediciones de latencia

---

## Reglas de Conectividad

### BR-08: Conexión Persistente
- **Regla**: El benchmark client usa una ÚNICA conexión TCP persistente para todas las iteraciones
- **No hay reconnect**: Si la conexión se pierde, el benchmark termina y reporta resultados parciales

### BR-09: Protocolo Síncrono
- **Regla**: El benchmark espera la respuesta antes de enviar el siguiente estímulo
- **Patrón**: Send → Wait → Receive → Measure → Repeat
- **No hay pipelining**: Cada iteración es un round-trip completo

---

## Extensibilidad del Payload (Documentación para Futuro)

El protocolo actual usa payload vacío (opción B). Aquí se documenta cómo implementar las alternativas para futuras versiones:

### Opción A: Payload "ping"/"pong" (4 bytes cada uno)

```
Formato del mensaje:
+--------+-------------------+
| Type   | Payload           |
| 1 byte | 4 bytes           |
+--------+-------------------+

Cambios requeridos:
1. protocol.go:
   - Encode: buf[0] = type, copy(buf[1:], payload) → return 1 + len(payload)
   - Decode: return data[0], data[1:], nil
   - Agregar constantes: StimulusPayload = []byte("ping"), ResponsePayload = []byte("pong")

2. handler.go:
   - Pre-allocar response = [5]byte{TypeResponse, 'p', 'o', 'n', 'g'}
   - conn.Write(response[:5])

3. benchmark/main.go:
   - Pre-allocar stimulus = [5]byte{TypeStimulus, 'p', 'i', 'n', 'g'}
   - readBuf = [5]byte{}

Impacto en latencia: +~1-5µs por la copia de 4 bytes adicionales
```

### Opción C: Payload Timestamp (8 bytes)

```
Formato del mensaje:
+--------+-------------------+
| Type   | Timestamp         |
| 1 byte | 8 bytes (int64)   |
+--------+-------------------+

Cambios requeridos:
1. protocol.go:
   - Encode: buf[0] = type, binary.LittleEndian.PutUint64(buf[1:], timestamp)
   - Decode: return data[0], binary.LittleEndian.Uint64(data[1:9]), nil
   - MessageSize = 9

2. handler.go:
   - Leer timestamp del stimulus
   - Incluir server timestamp en la response
   - Pre-allocar response = [9]byte{}

3. benchmark/main.go:
   - Incluir timestamp en el estímulo
   - Comparar timestamps del server para medir latencia server-side
   - readBuf = [9]byte{}

Impacto en latencia: +~5-10µs por serialización de int64 + copia de 8 bytes
Beneficio: Medición de latencia bidireccional (client→server y server→client)
```

---

## Propiedades Testables (PBT-01 Compliance)

| Componente | Propiedad | Categoría PBT | Enforced |
|---|---|---|---|
| `protocol.Encode/Decode` | Round-trip: `Decode(Encode(type)) == type` | Round-trip (PBT-02) | ✅ Sí |
| `protocol.Encode` | Invariante: output siempre tiene exactamente 1 byte | Invariant (PBT-03) | ✅ Sí |
| `protocol.Decode` | Invariante: tipo válido es 0x01 o 0x02 para inputs conocidos | Invariant (PBT-03) | ✅ Sí |
| `stats.Calculate` | Invariante: Min ≤ Median ≤ Max para cualquier input | Invariant (PBT-03) | ✅ Sí |
| `stats.Calculate` | Invariante: P50 ≤ P95 ≤ P99 para cualquier input | Invariant (PBT-03) | ✅ Sí |
| `stats.Calculate` | Invariante: Count == len(input) | Invariant (PBT-03) | ✅ Sí |
| `benchmarkLogger.Record/Flush` | Round-trip: datos registrados == datos escritos en archivo | Round-trip (PBT-02) | ✅ Sí |
