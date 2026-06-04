# Requerimientos del Sistema — Minimum Latency Challenge

## Intent Analysis Summary

- **User Request**: Diseñar y construir un sistema de estímulo-respuesta con latencia mínima (< 1ms)
- **Request Type**: New Project (Greenfield)
- **Scope Estimate**: Multiple Components (servidor, cliente benchmark, logging, documentación)
- **Complexity Estimate**: Moderate (optimización de rendimiento a nivel de microsegundos)
- **Depth**: Standard

---

## Requerimientos Funcionales

### RF-01: Servidor de Escucha Permanente
El sistema debe implementar un servidor que escuche permanentemente peticiones (estímulos) en un puerto TCP sobre localhost.

### RF-02: Respuesta a Estímulos
Al recibir un estímulo (mensaje binario), el sistema debe retornar una respuesta específica predefinida en formato binario.

### RF-03: Protocolo Binario Customizado
El protocolo de comunicación debe ser binario y customizado para máxima eficiencia, evitando overhead de serialización de texto (JSON, XML, etc.).

### RF-04: Medición de Latencia
El sistema debe medir el tiempo transcurrido (round-trip time) desde el momento en que se envía el estímulo hasta que se recibe la respuesta, usando relojes de alta resolución (nanosegundos).

### RF-05: Benchmark de 10,000 Iteraciones
El cliente de benchmark debe realizar 10,000 mediciones de latencia para generar estadísticas robustas incluyendo percentiles (p50, p95, p99).

### RF-06: Log de Trazabilidad
Se debe generar un archivo de log (.log) con nivel intermedio de detalle:
- Timestamp de envío, timestamp de recepción, latencia calculada por cada petición
- Estadísticas resumen al final: mínimo, máximo, promedio, mediana, percentiles p50, p95, p99

### RF-07: Documentación del Sistema
Se deben generar dos documentos Markdown:
1. **Documentación técnica**: Explicación de medición de latencia, justificación de herramientas/lenguaje, descripción de arquitectura
2. **Informe de resultados**: Comparación con el objetivo de < 1ms, análisis de resultados y optimizaciones

### RF-08: Reactor Pattern como Patrón Arquitectónico
El servidor debe implementar el **Reactor Pattern** como patrón de concurrencia principal:
- **Event Loop**: Un bucle de eventos single-threaded que demultiplexa y despacha eventos I/O
- **Event Demultiplexer**: Uso del mecanismo de notificación del OS (IOCP en Windows) para detectar eventos de red
- **Event Handlers**: Handlers registrados que procesan los estímulos cuando el demultiplexer notifica disponibilidad
- **Non-blocking I/O**: Todas las operaciones de red deben ser non-blocking para evitar bloqueo del event loop

Esto reemplaza el modelo estándar de Go (goroutine-per-connection) para eliminar overhead de scheduling y context switching.

---

## Requerimientos No Funcionales

### RNF-01: Latencia Objetivo
La latencia round-trip (envío del estímulo hasta recepción de la respuesta) debe ser menor a 1 milisegundo (1,000 microsegundos). El objetivo ideal es alcanzar el rango de microsegundos (< 100µs).

### RNF-02: Lenguaje de Implementación
El sistema será implementado en **Go** (Golang) por su balance entre rendimiento, simplicidad y concurrencia nativa via goroutines.

### RNF-03: Plataforma Objetivo
El sistema está optimizado para ejecución local en **Windows**. No se requiere portabilidad cross-platform.

### RNF-04: Comunicación TCP sobre Localhost
El mecanismo de comunicación será TCP sockets sobre `127.0.0.1` (localhost), proporcionando un canal de comunicación estándar y medible.

### RNF-05: Rendimiento y Eficiencia
- **Reactor Pattern**: Event loop single-threaded para eliminar goroutine scheduling overhead y context switching
- Uso de protocolo binario para minimizar overhead de serialización/deserialización
- Minimización de allocaciones de memoria en el hot path (zero-allocation en el hot path ideal)
- Uso de relojes de alta resolución (`time.Now()` con precisión de nanosegundos en Go)
- Reutilización de conexiones TCP (conexión persistente con el event loop)
- Pre-allocated buffers para evitar GC pressure
- Posible uso de frameworks como `gnet` o `evio` que implementan Reactor Pattern nativo en Go

### RNF-06: Fiabilidad del Servidor
El servidor debe manejar errores de conexión gracefully sin caerse, y poder atender múltiples peticiones secuenciales.

### RNF-07: Logging Estructurado
El sistema debe implementar logging estructurado con:
- Timestamp en formato ISO 8601
- Correlation/Request ID
- Log level
- Mensaje descriptivo
- Sin datos sensibles en los logs

---

## Restricciones

### Restricción de Tecnología
- **Lenguaje**: Go (confirmado por el usuario)
- **Package Manager**: No aplica (Go modules)
- **Plataforma**: Windows local

### Restricción de Negocio
- Es un proyecto experimental/académico enfocado en demostrar principios de arquitectura de software
- Contexto: Módulo 1 de Fundamentos de Arquitectura de Software (ISO/IEC 25010)

---

## Decisiones de Arquitectura

### ✅ Decididas

1. **Patrón de concurrencia**: **Reactor Pattern** — Event loop single-threaded con I/O non-blocking
   - Justificación: Elimina overhead de goroutine scheduling, minimiza context switching, latencia predecible
   - Trade-off (Primera Ley de Arquitectura): Se sacrifica la simplicidad idiomática de Go (goroutines) a cambio de rendimiento predecible y menor latencia
   - Implementación: Framework `gnet` (Reactor Pattern nativo para Go, basado en epoll/kqueue/IOCP) o implementación custom con `syscall`
2. **Patrón de comunicación**: TCP con conexión persistente sobre el event loop del Reactor
   - Justificación: Evita overhead de TCP handshake (3-way) en cada petición

### 🔲 Pendientes

1. **Estructura del protocolo binario**: Formato exacto del header y payload
2. **Framework vs Custom**: `gnet` vs `evio` vs implementación custom del Reactor
3. **Estrategia de warmup**: Pre-calentamiento de conexiones y event loop antes del benchmark
4. **Optimizaciones de Go runtime**: GOMAXPROCS, GC tuning, CPU affinity

---

## Extension Configuration

| Extension | Status | Enforcement Mode |
|---|---|---|
| Security Baseline | Enabled | Pragmático (proyecto experimental) |
| Property-Based Testing | Enabled | Partial (solo funciones puras y roundtrips de serialización) |

### Security Baseline — Aplicabilidad al Proyecto

| Rule | Aplicable | Justificación |
|---|---|---|
| SECURITY-01 | N/A | No hay data stores |
| SECURITY-02 | N/A | No hay load balancers ni API gateways |
| SECURITY-03 | Sí | Logging estructurado requerido |
| SECURITY-04 | N/A | No es una aplicación web |
| SECURITY-05 | Sí | Validación del protocolo binario |
| SECURITY-06 | N/A | No hay IAM |
| SECURITY-07 | N/A | Red local, no hay network policies |
| SECURITY-08 | N/A | No hay autenticación de usuarios |
| SECURITY-09 | Sí | Hardening básico, error handling |
| SECURITY-10 | Sí | Dependency pinning con go.sum |
| SECURITY-11 | Sí | Separación de concerns, rate limiting en benchmark |
| SECURITY-12 | N/A | No hay autenticación |
| SECURITY-13 | Sí | Validación de integridad en protocolo binario |
| SECURITY-14 | N/A | No requiere alerting para un benchmark local |
| SECURITY-15 | Sí | Error handling y fail-safe defaults |

### Property-Based Testing — Partial Enforcement

Reglas aplicables en modo parcial:
- **PBT-02**: Round-trip del protocolo binario (serialize/deserialize)
- **PBT-03**: Invariantes del protocolo (tamaño de mensajes, formato válido)
- **PBT-07**: Generadores de calidad para mensajes de protocolo
- **PBT-08**: Shrinking y reproducibilidad
- **PBT-09**: Framework selection (Go → `rapid`)

---

## Alineación con ISO/IEC 25010

| Característica | Subcaracterística | Relevancia | Prioridad |
|---|---|---|---|
| Eficiencia en el Desempeño | Comportamiento temporal | **CRÍTICA** — core del proyecto | ⭐⭐⭐ |
| Eficiencia en el Desempeño | Utilización de recursos | Alta — minimizar overhead | ⭐⭐ |
| Fiabilidad | Madurez | Media — servidor estable | ⭐⭐ |
| Fiabilidad | Tolerancia a fallos | Media — error handling | ⭐⭐ |
| Mantenibilidad | Modularidad | Alta — separación server/client/benchmark | ⭐⭐ |
| Mantenibilidad | Capacidad de prueba | Alta — benchmark y PBT | ⭐⭐ |
| Adecuación Funcional | Corrección funcional | Alta — respuesta correcta | ⭐⭐ |
