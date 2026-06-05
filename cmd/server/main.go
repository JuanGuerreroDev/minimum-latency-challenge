// Command server arranca el servidor TCP de mínima latencia basado en gnet.
//
// Usa un único event loop (Q1=A) con TCP_NODELAY para minimizar la latencia
// round-trip, escucha exclusivamente en localhost (PAT-SEC-01) y soporta
// graceful shutdown ante SIGINT/SIGTERM (Q2=A, PAT-RES-01).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/panjf2000/gnet/v2"

	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/logger"
	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/reactor"
)

func main() {
	port := flag.Int("port", 8080, "puerto TCP de escucha (localhost)")
	flag.Parse()

	log := logger.New(os.Stdout)

	handler := reactor.NewReactorHandler(log)

	// Bind exclusivo a localhost: mínima superficie de ataque (PAT-SEC-01).
	addr := fmt.Sprintf("tcp://127.0.0.1:%d", *port)

	// Captura de señales para graceful shutdown (PAT-RES-01).
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Info("shutdown signal received, stopping engine")
		if eng := handler.Engine(); eng != (gnet.Engine{}) {
			if err := eng.Stop(context.Background()); err != nil {
				log.Error("engine stop failed", err)
			}
		}
	}()

	log.Info("starting server", "addr", addr, "eventLoops", 1)

	// Single event loop, sin Nagle, para mínima latencia (PAT-PERF-01/02).
	err := gnet.Run(handler, addr,
		gnet.WithMulticore(false),
		gnet.WithNumEventLoop(1),
		gnet.WithReuseAddr(true),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
	)
	if err != nil {
		log.Error("server terminated", err)
		os.Exit(1)
	}
}
