# Component Methods — Minimum Latency Challenge

> **Nota**: Este documento define las firmas de métodos de alto nivel. La lógica de negocio detallada se definirá en la fase de Functional Design (CONSTRUCTION).

---

## 1. Server (`cmd/server`)

### `func main()`
- **Propósito**: Entry point del servidor
- **Input**: Ninguno (usa flags/env para configuración)
- **Output**: Ninguno (corre indefinidamente)
- **Notas**: Inicializa reactor handler, arranca gnet event loop

---

## 2. Benchmark Client (`cmd/benchmark`)

### `func main()`
- **Propósito**: Entry point del benchmark
- **Input**: Ninguno (usa flags para host:port, iterations)
- **Output**: Ninguno (escribe a stdout y archivo .log)
- **Notas**: Conecta al server, ejecuta benchmark, reporta resultados

### `func runBenchmark(conn net.Conn, iterations int) []time.Duration`
- **Propósito**: Ejecuta el loop principal de benchmark
- **Input**: `conn` (conexión TCP), `iterations` (número de iteraciones)
- **Output**: Slice de duraciones medidas
- **Notas**: Hot path — performance critical, zero-allocation ideal

---

## 3. Reactor Handler (`internal/reactor`)

### `type ReactorHandler struct`
- **Campos**: `logger`, response buffer pre-allocated

### `func (h *ReactorHandler) OnBoot(eng gnet.Engine) gnet.Action`
- **Propósito**: Callback cuando el event loop arranca
- **Input**: `eng` (engine de gnet)
- **Output**: `gnet.Action` (None para continuar)

### `func (h *ReactorHandler) OnOpen(c gnet.Conn) ([]byte, gnet.Action)`
- **Propósito**: Callback cuando una nueva conexión se abre
- **Input**: `c` (conexión gnet)
- **Output**: Datos iniciales opcionales, action

### `func (h *ReactorHandler) OnClose(c gnet.Conn, err error) gnet.Action`
- **Propósito**: Callback cuando una conexión se cierra
- **Input**: `c` (conexión), `err` (error si hubo)
- **Output**: `gnet.Action`

### `func (h *ReactorHandler) OnTraffic(c gnet.Conn) gnet.Action`
- **Propósito**: **HOT PATH** — Callback cuando hay datos disponibles para leer
- **Input**: `c` (conexión gnet con buffer de lectura)
- **Output**: `gnet.Action`
- **Notas**: Decodifica estímulo, codifica respuesta, escribe directo. Zero-allocation en este método es CRÍTICO.

### `func NewReactorHandler(logger *logger.Logger) *ReactorHandler`
- **Propósito**: Constructor del handler
- **Input**: `logger` (logger estructurado)
- **Output**: Handler configurado con buffers pre-allocated

---

## 4. Protocol (`internal/protocol`)

### `func Encode(msgType byte, payload []byte, buf []byte) int`
- **Propósito**: Codifica un mensaje al formato binario en un buffer provisto
- **Input**: `msgType` (tipo de mensaje), `payload` (datos), `buf` (buffer destino pre-allocated)
- **Output**: Número de bytes escritos
- **Notas**: Zero-allocation — usa buffer pre-allocated del caller

### `func Decode(data []byte) (msgType byte, payload []byte, err error)`
- **Propósito**: Decodifica un mensaje desde formato binario
- **Input**: `data` (bytes crudos)
- **Output**: Tipo de mensaje, payload (slice del input, no copia), error
- **Notas**: Zero-allocation — retorna slice del input original

### Constantes
```go
const (
    TypeStimulus byte = 0x01
    TypeResponse byte = 0x02
    HeaderSize   int  = 1
)
```

---

## 5. Logger (`internal/logger`)

### `type Logger struct`
- **Campos**: Output writer, level, buffer de mediciones

### `func New(output io.Writer) *Logger`
- **Propósito**: Crea un nuevo logger estructurado
- **Input**: `output` (writer para logs del sistema)
- **Output**: Logger configurado

### `func (l *Logger) Info(msg string, fields ...Field)`
- **Propósito**: Log de nivel INFO con campos estructurados
- **Input**: Mensaje y campos clave-valor
- **Output**: Ninguno

### `func (l *Logger) Error(msg string, err error, fields ...Field)`
- **Propósito**: Log de nivel ERROR
- **Input**: Mensaje, error, campos adicionales
- **Output**: Ninguno

### `type BenchmarkLogger struct`
- **Campos**: Buffer de registros pre-allocated (10,000 slots)

### `func NewBenchmarkLogger(capacity int) *BenchmarkLogger`
- **Propósito**: Crea logger de benchmark con buffer pre-allocated
- **Input**: Capacidad del buffer
- **Output**: BenchmarkLogger listo

### `func (bl *BenchmarkLogger) Record(sendTime, recvTime time.Time, latency time.Duration)`
- **Propósito**: Registra una medición en el buffer (zero-allocation)
- **Input**: Timestamps y latencia
- **Output**: Ninguno

### `func (bl *BenchmarkLogger) FlushToFile(filename string, stats *stats.Stats) error`
- **Propósito**: Escribe todos los registros + estadísticas al archivo .log
- **Input**: Nombre del archivo, estadísticas calculadas
- **Output**: Error si falla la escritura

---

## 6. Stats (`internal/stats`)

### `type Stats struct`
- **Campos**: Min, Max, Avg, Median, P50, P95, P99, Count, Total

### `func Calculate(durations []time.Duration) *Stats`
- **Propósito**: Calcula todas las estadísticas a partir de las mediciones
- **Input**: Slice de duraciones
- **Output**: Estadísticas completas

### `func (s *Stats) String() string`
- **Propósito**: Formato legible de las estadísticas
- **Input**: Ninguno
- **Output**: String formateado con todas las métricas

### `func (s *Stats) Report() string`
- **Propósito**: Reporte detallado para stdout con formato tabular
- **Input**: Ninguno
- **Output**: Reporte completo
