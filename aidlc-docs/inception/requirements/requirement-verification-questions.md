# Preguntas de Verificación de Requerimientos

Por favor responde las siguientes preguntas para clarificar los requerimientos del sistema de latencia mínima. Llena la etiqueta `[Answer]:` con la letra de tu elección.

## Question 1

¿Qué mecanismo de comunicación prefieres para el estímulo-respuesta?

A) IPC (Inter-Process Communication) en la misma máquina — Unix Domain Sockets, pipes, shared memory
B) Comunicación en red local — TCP/UDP sockets sobre localhost
C) Ambos — implementar IPC como principal y TCP/UDP como alternativa para comparar
D) Comunicación intra-proceso — canales/goroutines dentro del mismo proceso
X) Other (please describe after [Answer]: tag below)

[Answer]: B

## Question 2

¿Cuál es el formato preferido del estímulo (mensaje de entrada)?

A) String simple — un mensaje de texto plano (ej: "ping")
B) JSON — un objeto JSON estructurado (ej: {"type": "stimulus", "data": "..."})
C) Binario — protocolo binario customizado para máxima eficiencia
D) Protobuf/FlatBuffers — serialización binaria con schema
X) Other (please describe after [Answer]: tag below)

[Answer]: C

## Question 3

¿Cuántas mediciones de latencia se deben realizar para generar estadísticas confiables?

A) 100 iteraciones — rápido, suficiente para una demo
B) 1,000 iteraciones — buen balance entre precisión y velocidad
C) 10,000 iteraciones — estadísticas más robustas con percentiles (p50, p95, p99)
D) 100,000 iteraciones — análisis exhaustivo de distribución de latencia
X) Other (please describe after [Answer]: tag below)

[Answer]: C

## Question 4

¿Qué nivel de detalle esperas en el log de trazabilidad?

A) Básico — solo timestamp de envío, timestamp de recepción y latencia calculada por cada petición
B) Intermedio — lo anterior + estadísticas resumen al final (min, max, promedio, percentiles)
C) Detallado — lo anterior + información del sistema (CPU, memoria, OS) y parámetros de configuración
X) Other (please describe after [Answer]: tag below)

[Answer]: B

## Question 5

¿El sistema debe poder ejecutarse solo en tu máquina local (Windows) o debe ser portable a otros entornos?

A) Solo local Windows — optimizado para mi máquina
B) Cross-platform — debe compilar y correr en Windows, Linux y macOS
C) Principalmente Linux — orientado a despliegue en servidores Linux pero que compile en Windows
X) Other (please describe after [Answer]: tag below)

[Answer]: A

## Question 6

¿Confirmas Go como el lenguaje principal para la implementación?

A) Sí, Go — buen balance entre rendimiento, simplicidad y concurrencia nativa
B) Rust — máximo rendimiento y control de memoria, pero mayor complejidad
C) C/C++ — rendimiento nativo puro, pero más propenso a errores
D) Quiero una comparación — implementar en Go y en otro lenguaje para comparar latencias
X) Other (please describe after [Answer]: tag below)

[Answer]: A

## Question 7: Security Extensions

¿Deben aplicarse las reglas de la extensión de seguridad (Security Baseline) en este proyecto?

A) Sí — aplicar todas las reglas de SEGURIDAD como restricciones bloqueantes (recomendado para aplicaciones de producción)
B) No — omitir todas las reglas de SEGURIDAD (adecuado para PoCs, prototipos y proyectos experimentales)
X) Other (please describe after [Answer]: tag below)

[Answer]: A, pero permitir la opción de desactivarlas, ya que este es un proyecto experimental.

## Question 8: Property-Based Testing Extension

¿Deben aplicarse las reglas de Property-Based Testing (PBT) en este proyecto?

A) Sí — aplicar todas las reglas de PBT como restricciones bloqueantes (recomendado para proyectos con lógica de negocio, transformaciones de datos o componentes con estado)
B) Parcial — aplicar reglas de PBT solo para funciones puras y roundtrips de serialización
C) No — omitir todas las reglas de PBT (adecuado para aplicaciones CRUD simples, proyectos solo de UI, o capas de integración delgadas)
X) Other (please describe after [Answer]: tag below)

[Answer]: B
