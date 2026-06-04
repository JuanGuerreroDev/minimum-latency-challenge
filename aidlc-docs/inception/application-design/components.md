# Componentes del Sistema — Minimum Latency Challenge

## Visión General

El sistema se compone de 6 componentes organizados en una estructura package-based de Go, implementando el Reactor Pattern mediante el framework `gnet` para comunicación TCP de ultra-baja latencia.

---

## 1. Server (`cmd/server`)

**Propósito**: Entry point del servidor TCP basado en gnet que escucha permanentemente estímulos.

**Responsabilidades**:
- Inicializar el event loop de gnet en el puerto TCP configurado
- Configurar el logger estructurado
- Manejar señales del OS para graceful shutdown
- Imprimir información de startup (puerto, configuración)

**Tipo**: Ejecutable (`main` package)

---

## 2. Benchmark Client (`cmd/benchmark`)

**Propósito**: Entry point del cliente de benchmark que mide la latencia round-trip.

**Responsabilidades**:
- Establecer conexión TCP persistente al servidor
- Ejecutar 10,000 peticiones en modo single-shot
- Medir latencia round-trip con relojes de alta resolución (nanosegundos)
- Acumular resultados en buffer de memoria
- Calcular estadísticas de latencia (min, max, avg, p50, p95, p99)
- Escribir archivo de log (.log) con trazabilidad al finalizar
- Reportar resultados por stdout

**Tipo**: Ejecutable (`main` package)

---

## 3. Reactor Handler (`internal/reactor`)

**Propósito**: Implementación del Event Handler de gnet que procesa los estímulos según el Reactor Pattern.

**Responsabilidades**:
- Implementar la interfaz `gnet.EventHandler` de gnet
- Manejar eventos del ciclo de vida: `OnBoot`, `OnOpen`, `OnClose`, `OnShutdown`
- Procesar datos entrantes en `OnTraffic` (evento de I/O)
- Decodificar el estímulo binario usando el protocolo
- Codificar y enviar la respuesta binaria
- Manejar errores de conexión sin crashear el event loop

**Tipo**: Package interno (`internal/reactor`)

---

## 4. Protocol (`internal/protocol`)

**Propósito**: Definición y serialización/deserialización del protocolo binario ultra-minimal.

**Responsabilidades**:
- Definir tipos de mensaje (stimulus, response)
- Codificar mensajes a formato binario (1 byte tipo + N bytes payload)
- Decodificar mensajes desde formato binario
- Definir constantes del protocolo (tipos, tamaños máximos)
- Validar integridad básica de los mensajes recibidos

**Tipo**: Package interno (`internal/protocol`)

**Formato del Protocolo**:
```
+--------+-------------------+
| Type   | Payload           |
| 1 byte | Variable (0-255)  |
+--------+-------------------+
```

- **Type** (1 byte): `0x01` = Stimulus, `0x02` = Response
- **Payload**: Datos del mensaje (para el benchmark: string fijo "ping" / "pong")

---

## 5. Logger (`internal/logger`)

**Propósito**: Logging estructurado del sistema y escritura buffered del archivo de trazabilidad.

**Responsabilidades**:
- Proveer logger estructurado con timestamp ISO 8601, level, y message
- Acumular registros de latencia en buffer de memoria durante el benchmark
- Flush del buffer al archivo .log al finalizar el benchmark
- Formatear cada registro con: timestamp envío, timestamp recepción, latencia
- Generar sección de estadísticas resumen al final del log

**Tipo**: Package interno (`internal/logger`)

---

## 6. Stats (`internal/stats`)

**Propósito**: Cálculo de estadísticas de latencia a partir de las mediciones del benchmark.

**Responsabilidades**:
- Almacenar slice pre-allocated de mediciones de latencia (10,000 slots)
- Calcular mínimo, máximo, promedio, mediana
- Calcular percentiles p50, p95, p99
- Generar reporte formateado de estadísticas
- Proveer datos para el archivo de log y stdout

**Tipo**: Package interno (`internal/stats`)
