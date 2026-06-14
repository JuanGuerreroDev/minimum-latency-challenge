// Command benchmark mide la latencia round-trip contra el servidor de mínima
// latencia: envía un estímulo (0x01) y espera la respuesta (0x02), repitiendo
// N iteraciones sobre una única conexión TCP persistente (BR-08, BR-09).
//
// El loop de medición es síncrono (send → wait → recv → measure), no hace I/O
// de disco ni logging durante la ejecución (BR-06, BR-07), y aplica
// skip & continue ante errores individuales (BR-04, PAT-RES-02).
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/logger"
	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/protocol"
	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/stats"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host del servidor")
	port := flag.Int("port", 8080, "puerto del servidor")
	iterations := flag.Int("iterations", 10000, "número de iteraciones del benchmark")
	logFile := flag.String("log", "benchmark.log", "archivo de trazabilidad de salida")
	flag.Parse()

	log := logger.New(os.Stdout)
	addr := net.JoinHostPort(*host, strconv.Itoa(*port))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error("no se pudo conectar al servidor", err, "addr", addr)
		os.Exit(1)
	}
	defer conn.Close()

	// Deshabilitar Nagle en el cliente para minimizar latencia (PAT-PERF-02).
	if tcp, ok := conn.(*net.TCPConn); ok {
		if err := tcp.SetNoDelay(true); err != nil {
			log.Error("no se pudo deshabilitar Nagle", err)
		}
	}

	log.Info("benchmark started", "addr", addr, "iterations", *iterations)

	benchLog := logger.NewLatencyRecorder(*iterations)
	durations, errCount := runBenchmark(conn, *iterations, benchLog)

	s := stats.Calculate(durations)

	// Reporte solo al final (BR-06).
	fmt.Print(s.Report())

	// Flush del archivo de trazabilidad post-benchmark (BR-07, NFR-LOG-02).
	if err := benchLog.FlushToFile(*logFile, s); err != nil {
		log.Error("no se pudo escribir el log de benchmark", err, "file", *logFile)
	} else {
		fmt.Printf("Trace log escrito en %s\n", *logFile)
	}

	if errCount > 0 {
		fmt.Printf("WARNING: %d iteraciones fallaron (excluidas de las estadísticas)\n", errCount)
	}
}

// runBenchmark ejecuta el loop de medición. HOT PATH: usa buffers
// pre-allocated y un slice de duraciones con capacidad reservada para evitar
// asignaciones durante la medición (PAT-PERF-03).
//
// Retorna las duraciones de las iteraciones exitosas y el número de errores.
func runBenchmark(conn net.Conn, iterations int, benchLog *logger.LatencyRecorder) ([]time.Duration, int) {
	var (
		stimulus  [protocol.MessageSize]byte
		readBuf   [protocol.MessageSize]byte
		durations = make([]time.Duration, 0, iterations)
		errCount  int
	)
	protocol.Encode(protocol.TypeStimulus, stimulus[:])

	for i := 0; i < iterations; i++ {
		sendTime := time.Now()

		if _, err := conn.Write(stimulus[:]); err != nil {
			errCount++
			continue // skip & continue (BR-04)
		}
		if _, err := conn.Read(readBuf[:]); err != nil {
			errCount++
			continue // skip & continue (BR-04)
		}

		recvTime := time.Now()
		latency := recvTime.Sub(sendTime)
		durations = append(durations, latency)
		benchLog.Record(sendTime, recvTime, latency)
	}

	return durations, errCount
}
