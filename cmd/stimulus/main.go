// Command stimulus provee un cliente TCP interactivo que permite enviar
// estímulos (0x01) al servidor de mínima latencia en cualquier momento.
//
// El usuario presiona Enter para enviar un estímulo y ver la latencia
// round-trip. La sesión permanece activa con una conexión TCP persistente
// hasta que se recibe SIGINT (Ctrl+C) o EOF en stdin.
//
// Cumple RF-09 (estímulo interactivo) y RF-10 (trazabilidad con log).
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/logger"
	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/protocol"
	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/stats"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host del servidor")
	port := flag.Int("port", 8080, "puerto del servidor")
	logFile := flag.String("log", "stimulus.log", "archivo de trazabilidad de salida")
	flag.Parse()

	log := logger.New(os.Stdout)
	addr := net.JoinHostPort(*host, strconv.Itoa(*port))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error("no se pudo conectar al servidor", err, "addr", addr)
		os.Exit(1)
	}
	defer conn.Close()

	// Deshabilitar Nagle para minimizar latencia (PAT-PERF-02).
	if tcp, ok := conn.(*net.TCPConn); ok {
		if err := tcp.SetNoDelay(true); err != nil {
			log.Error("no se pudo deshabilitar Nagle", err)
		}
	}

	// Graceful shutdown: capturar SIGINT/SIGTERM (PAT-RES-01).
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Buffers pre-allocated para el hot path (zero-allocation, PAT-PERF-03).
	var stimulus [protocol.MessageSize]byte
	var readBuf [protocol.MessageSize]byte
	protocol.Encode(protocol.TypeStimulus, stimulus[:])

	// LatencyRecorder para acumular trazas (capacidad inicial 64, crece si necesita).
	benchLog := logger.NewLatencyRecorder(64)
	// Tracking local de duraciones para estadísticas al final.
	durations := make([]time.Duration, 0, 64)

	fmt.Println("=== Stimulus Client (Interactive Mode) ===")
	fmt.Printf("Conectado a %s\n", addr)
	fmt.Println("Presiona Enter para enviar un estímulo. Ctrl+C para salir.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	// Canal para detectar EOF en stdin.
	done := make(chan struct{})
	go func() {
		for scanner.Scan() {
			// Cada Enter dispara un estímulo.
			sendTime := time.Now()

			if _, werr := conn.Write(stimulus[:]); werr != nil {
				log.Error("error al enviar estímulo", werr)
				continue
			}

			if _, rerr := conn.Read(readBuf[:]); rerr != nil {
				log.Error("error al leer respuesta", rerr)
				continue
			}

			recvTime := time.Now()
			latency := recvTime.Sub(sendTime)

			fmt.Printf("  → Respuesta recibida | Latency: %v\n", latency)
			benchLog.Record(sendTime, recvTime, latency)
			durations = append(durations, latency)
		}
		close(done)
	}()

	// Esperar señal de shutdown o EOF en stdin.
	select {
	case <-ctx.Done():
		fmt.Println("\nShutdown signal received.")
	case <-done:
		fmt.Println("\nEOF en stdin.")
	}

	// Flush del log de trazabilidad (RF-10).
	count := benchLog.Count()
	if count > 0 {
		var s *stats.Stats
		if count >= 2 {
			s = stats.Calculate(durations)
			fmt.Print(s.Report())
		}
		if err := benchLog.FlushToFile(*logFile, s); err != nil {
			log.Error("no se pudo escribir el log", err, "file", *logFile)
		} else {
			fmt.Printf("Trace log escrito en %s (%d estímulos)\n", *logFile, count)
		}
	} else {
		fmt.Println("No se enviaron estímulos. Sin log generado.")
	}
}

