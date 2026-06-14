---
status: accepted
date: 2026-06-13
decision-makers: Juan Guerrero, Fabian Capote, Estefania Marquez, Pedro Tabares 
consulted: AI-DLC
informed: Equipo de Diplomado en Arquitectura de Software
---

# Fundamentos de diseño: latencia, protocolo y modularidad

## Contexto y Planteamiento del Problema

Se necesita un sistema de estímulo-respuesta TCP con latencia de ida y vuelta p99 < 1ms en localhost. El sistema debe demostrar principios de arquitectura de software orientados a rendimiento, permitir medición precisa de latencia vía benchmark automatizado, y además permitir el envío de estímulos individuales en cualquier momento de forma interactiva.

¿Qué patrón arquitectónico, protocolo de comunicación y estructura de binarios permiten alcanzar latencia sub-milisegundo manteniendo modularidad y claridad de responsabilidades?

## Factores de Decisión

- Latencia de ida y vuelta p99 < 1ms como objetivo principal (RNF dominante)
- Cero asignaciones de memoria en la ruta crítica para eliminar presión sobre el recolector de basura
- Separación de responsabilidades (SRP) entre componentes
- Simplicidad de implementación dentro de las restricciones de rendimiento
- Capacidad de medir latencia con precisión (nanosegundos)
- Envío de estímulos tanto en modo batch (benchmark) como ad-hoc (interactivo)
- Contexto académico: demostrar trade-offs arquitectónicos (ISO/IEC 25010)

## Opciones Consideradas

### Patrón arquitectónico
- **Reactor Pattern (gnet)** — Bucle de eventos de un solo hilo con E/S no bloqueante
- **Goroutine por conexión** — Modelo idiomático de Go con una goroutine por conexión
- **Syscall custom** — Implementación directa sobre IOCP/epoll sin framework

### Protocolo
- **Binario ultra-minimal (1 byte)** — Solo tipo de mensaje, sin payload
- **Binario con payload** — Tipo + datos (ej: "ping"/"pong")
- **Texto (JSON/HTTP)** — Protocolo estándar sobre TCP

### Estructura de binarios
- **Binarios separados** — `server`, `benchmark`, `stimulus` independientes
- **Binario unificado con subcomandos** — Un solo ejecutable con modos
- **Flag en benchmark** — Agregar modo interactivo al benchmark existente

## Resultado de la Decisión

### 1. Reactor Pattern (gnet) sobre goroutine por conexión

Elegido porque elimina el overhead de planificación de goroutines y cambio de contexto. Con un solo bucle de eventos (`NumEventLoop=1`, `Multicore=false`) se maximiza la localidad de caché para el caso de uso principal. gnet provee Reactor Pattern nativo para Go con soporte IOCP en Windows, probado en producción.

### 2. Protocolo binario ultra-minimal de 1 byte

Elegido porque elimina completamente el overhead de serialización. Un mensaje = 1 byte (`0x01` = Estímulo, `0x02` = Respuesta). Las funciones `Encode`/`Decode` son puras y no asignan memoria. El protocolo es tan simple que no hay ramificación compleja en la ruta crítica.

### 3. Binarios separados (server + benchmark + stimulus)

Elegido porque respeta SRP: cada binario tiene una sola responsabilidad. El benchmark mide rendimiento en bucle cerrado; el stimulus permite envío interactivo ad-hoc. Comparten paquetes internos (`protocol`, `logger`) sin acoplamiento entre sí.

### 4. LatencyRecorder como abstracción genérica (renombramiento de BenchmarkLogger)

Elegido porque el nombre original (`BenchmarkLogger`) acoplaba semánticamente al benchmark una estructura cuya responsabilidad real es "acumular y persistir mediciones de latencia". Con el renombramiento a `LatencyRecorder` se refleja SRP y se habilita reutilización limpia desde cualquier cliente.

### Consecuencias

- Positivo: p99 = 646µs (benchmark) y 555µs (prueba corta) — objetivo < 1ms alcanzado
- Positivo: 0 asignaciones/operación en ruta crítica verificado con benchmarks de Go
- Positivo: modularidad: 6 paquetes con responsabilidad única y dependencias claras
- Positivo: el protocolo es extensible (bytes 0x03+ disponibles) sin romper compatibilidad
- Positivo: dos modos de cliente cubren tanto medición rigurosa como uso ad-hoc
- Negativo: se sacrifica la idiomática de Go (goroutines) a favor de rendimiento
- Negativo: el modo interactivo tiene fluctuaciones inherentes por planificación del SO (p99 > 1ms en uso humano)
- Neutral: el bucle de eventos único es óptimo para 1 conexión pero necesitaría ajuste si se escala a cientos

## Plan de Implementación

- **Rutas afectadas**:
  - `cmd/server/main.go` — Punto de entrada del servidor con gnet
  - `cmd/benchmark/main.go` — Cliente de medición en bucle cerrado
  - `cmd/stimulus/main.go` — Cliente interactivo (Enter → estímulo)
  - `internal/reactor/handler.go` — Manejador de eventos (OnTraffic = ruta crítica)
  - `internal/protocol/protocol.go` — Encode/Decode del protocolo binario
  - `internal/logger/latency_recorder.go` — Acumulador buffered de mediciones
  - `internal/stats/stats.go` — Cálculo de percentiles
- **Dependencias**: `github.com/panjf2000/gnet/v2` (Reactor Pattern), `pgregory.net/rapid` (pruebas basadas en propiedades)
- **Patrones a seguir**:
  - Buffers pre-asignados fuera del bucle para cero asignaciones
  - Funciones puras en `protocol` y `stats` (sin efectos secundarios)
  - Logging solo en rutas frías; nunca en la ruta crítica exitosa
  - Cierre ordenado con `signal.NotifyContext` + vaciado de logs
  - TCP_NODELAY en todos los endpoints (anti-Nagle)
- **Patrones a evitar**:
  - No asignar memoria en `OnTraffic` ni en el bucle de envío/recepción
  - No hacer E/S de disco durante mediciones de latencia
  - No usar goroutine por conexión para el servidor
  - No usar protocolos de texto (JSON/HTTP) en el canal de comunicación
- **Configuración**: Flags de línea de comandos (`--host`, `--port`, `--log`, `--iterations`)

### Verificación

- [x] `go build ./...` compila los 3 binarios sin errores
- [x] `go test ./...` — 12 pruebas pasan (incluyendo pruebas basadas en propiedades)
- [x] `go vet ./...` — sin advertencias
- [x] Benchmark de 10,000 iteraciones: p99 < 1ms
- [x] Stimulus interactivo: envía 0x01, recibe 0x02, muestra latencia
- [x] Log de trazabilidad generado correctamente (formato consistente con benchmark)
- [x] Cierre ordenado funciona (SIGINT + vaciado del log)
- [x] Cero referencias a `BenchmarkLogger` en el código fuente (renombramiento completo)

## Pros y Contras de las Opciones

### Reactor Pattern (gnet)

- Positivo: latencia predecible sin overhead de planificación de goroutines
- Positivo: IOCP nativo en Windows sin código custom
- Positivo: framework probado en producción (ecosistema con >100k estrellas en GitHub)
- Neutral: requiere aprender la API de callbacks de gnet
- Negativo: no es idiomático en Go (sacrifica simplicidad de goroutines)
- Negativo: bucle de eventos único es subóptimo si se necesitan muchos núcleos

### Goroutine por conexión

- Positivo: idiomático en Go, fácil de entender
- Positivo: escala naturalmente a múltiples núcleos
- Negativo: overhead de planificación (~1-5µs por cambio de contexto)
- Negativo: presión sobre el recolector de basura si se asignan buffers por goroutine
- Negativo: latencia menos predecible bajo carga

### Protocolo binario de 1 byte

- Positivo: cero overhead de serialización
- Positivo: funciones Encode/Decode triviales y sin asignaciones
- Positivo: extensible sin romper compatibilidad (bytes 0x03-0xFF disponibles)
- Neutral: no es legible por humanos (necesita herramientas para depuración)
- Negativo: no soporta payload variable sin rediseño

### Protocolo con payload (tipo + datos)

- Positivo: soporta mensajes con contenido variable
- Neutral: overhead marginal (1-2 bytes extra de prefijo de longitud)
- Negativo: requiere análisis más complejo en la ruta crítica
- Negativo: introduce asignaciones para payload dinámico

## Información Adicional

### Condiciones de revisión

Reconsiderar estas decisiones si:
- Se necesita ejecutar en Linux para una demo o comparación IOCP vs epoll
- El objetivo de latencia baja de 1ms a <100µs p99 consistente
- Se requiere soportar cientos de conexiones simultáneas bajo carga
- El proyecto evoluciona de experimental a producción

### Referencias

- Documentación de arquitectura: `docs/system-documentation.md`
- Resultados de benchmark: `docs/results-report.md`
- gnet: https://github.com/panjf2000/gnet
- ISO/IEC 25010 (Atributos de calidad): marco de evaluación del proyecto
