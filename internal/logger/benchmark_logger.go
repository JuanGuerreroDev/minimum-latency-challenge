package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/stats"
)

// BenchmarkLogger acumula mediciones de latencia en memoria durante el
// benchmark y las vuelca a un archivo .log al finalizar.
//
// Usa tres slices paralelos pre-allocated en vez de un slice de structs para
// mejorar la localidad de caché del write path (Record), y para garantizar
// zero-allocation por iteración: Record solo escribe en posiciones ya
// reservadas (PAT-OBS-02, PAT-PERF-03, BR-07).
type BenchmarkLogger struct {
	sendTimes []time.Time
	recvTimes []time.Time
	latencies []time.Duration
	count     int
}

// NewBenchmarkLogger crea un BenchmarkLogger con los buffers pre-allocated a la
// capacidad indicada (típicamente el número de iteraciones del benchmark).
func NewBenchmarkLogger(capacity int) *BenchmarkLogger {
	if capacity < 0 {
		capacity = 0
	}
	return &BenchmarkLogger{
		sendTimes: make([]time.Time, capacity),
		recvTimes: make([]time.Time, capacity),
		latencies: make([]time.Duration, capacity),
	}
}

// Record almacena una medición individual. Es zero-allocation: escribe en slots
// pre-allocated y no realiza I/O (cero impacto en las mediciones de latencia).
//
// Si se supera la capacidad reservada, los slices crecen (append), lo cual solo
// ocurriría si se registran más mediciones que la capacidad inicial.
func (bl *BenchmarkLogger) Record(sendTime, recvTime time.Time, latency time.Duration) {
	if bl.count < len(bl.latencies) {
		bl.sendTimes[bl.count] = sendTime
		bl.recvTimes[bl.count] = recvTime
		bl.latencies[bl.count] = latency
	} else {
		bl.sendTimes = append(bl.sendTimes, sendTime)
		bl.recvTimes = append(bl.recvTimes, recvTime)
		bl.latencies = append(bl.latencies, latency)
	}
	bl.count++
}

// FlushToFile escribe todas las mediciones registradas más la sección de
// estadísticas al archivo indicado, usando escritura buffered. Se invoca
// post-benchmark para no afectar las mediciones (BR-07, NFR-LOG-02).
func (bl *BenchmarkLogger) FlushToFile(filename string, s *stats.Stats) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("logger: no se pudo crear %q: %w", filename, err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	fmt.Fprintln(w, "# Benchmark Trace Log — Minimum Latency System")
	fmt.Fprintf(w, "# Iterations recorded: %d\n", bl.count)
	fmt.Fprintln(w, "#")
	fmt.Fprintln(w, "# SendTime | RecvTime | Latency")
	fmt.Fprintln(w, "# -----------------------------------------------")

	for i := 0; i < bl.count; i++ {
		fmt.Fprintf(w, "%s | %s | %v\n",
			bl.sendTimes[i].Format(time.RFC3339Nano),
			bl.recvTimes[i].Format(time.RFC3339Nano),
			bl.latencies[i],
		)
	}

	fmt.Fprintln(w, "# -----------------------------------------------")
	fmt.Fprintln(w, "# Summary Statistics")
	if s != nil {
		fmt.Fprintf(w, "#   Count : %d\n", s.Count)
		fmt.Fprintf(w, "#   Min   : %v\n", s.Min)
		fmt.Fprintf(w, "#   Max   : %v\n", s.Max)
		fmt.Fprintf(w, "#   Avg   : %v\n", s.Avg)
		fmt.Fprintf(w, "#   p50   : %v\n", s.P50)
		fmt.Fprintf(w, "#   p95   : %v\n", s.P95)
		fmt.Fprintf(w, "#   p99   : %v\n", s.P99)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("logger: error al hacer flush de %q: %w", filename, err)
	}
	return nil
}

// Count retorna el número de mediciones registradas hasta el momento.
func (bl *BenchmarkLogger) Count() int {
	return bl.count
}
