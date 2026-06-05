// Package reactor implementa el Event Handler del Reactor Pattern sobre gnet.
//
// OnTraffic es el HOT PATH del sistema: por cada estímulo recibido responde con
// un único byte desde un buffer pre-allocated, sin asignar memoria
// (PAT-PERF-01/02/03). El resto de callbacks del ciclo de vida son cold path y
// solo registran eventos.
package reactor

import (
	"github.com/panjf2000/gnet/v2"

	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/logger"
	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/protocol"
)

// ReactorHandler implementa gnet.EventHandler para el servidor de mínima
// latencia. Embebe gnet.BuiltinEventEngine para heredar implementaciones por
// defecto de los callbacks que no se sobreescriben.
type ReactorHandler struct {
	gnet.BuiltinEventEngine

	log      *logger.Logger
	eng      gnet.Engine
	response [protocol.MessageSize]byte // buffer pre-allocated, reutilizado (0 allocs)
}

// NewReactorHandler crea el handler con su buffer de respuesta pre-allocated.
func NewReactorHandler(log *logger.Logger) *ReactorHandler {
	h := &ReactorHandler{log: log}
	// El tipo de respuesta es constante; se fija una sola vez.
	h.response[0] = protocol.TypeResponse
	return h
}

// Engine expone el engine de gnet capturado en OnBoot, para que el caller pueda
// ejecutar un shutdown limpio (PAT-RES-01).
func (h *ReactorHandler) Engine() gnet.Engine {
	return h.eng
}

// OnBoot se invoca una vez cuando el event loop arranca (cold path).
func (h *ReactorHandler) OnBoot(eng gnet.Engine) gnet.Action {
	h.eng = eng
	h.log.Info("server started")
	return gnet.None
}

// OnShutdown se invoca cuando el engine se detiene (cold path).
func (h *ReactorHandler) OnShutdown(_ gnet.Engine) {
	h.log.Info("server shutdown complete")
}

// OnOpen se invoca al abrir una conexión (cold path). No envía datos iniciales.
func (h *ReactorHandler) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	h.log.Info("connection opened", "remote", c.RemoteAddr().String())
	return nil, gnet.None
}

// OnClose se invoca al cerrar una conexión (cold path).
func (h *ReactorHandler) OnClose(c gnet.Conn, err error) gnet.Action {
	if err != nil {
		h.log.Error("connection closed with error", err, "remote", c.RemoteAddr().String())
	} else {
		h.log.Info("connection closed", "remote", c.RemoteAddr().String())
	}
	return gnet.None
}

// OnTraffic es el HOT PATH: se ejecuta por cada lote de datos disponibles.
//
// Invariantes (PAT-PERF-03):
//   - Zero allocations: usa c.Next (vista al buffer interno de gnet) y el buffer
//     de respuesta pre-allocated del handler.
//   - Branching mínimo: un único if/else.
//   - Fail closed: ante tipo inválido no responde, solo loggea (BR-01, PAT-SEC-02).
func (h *ReactorHandler) OnTraffic(c gnet.Conn) gnet.Action {
	buf, err := c.Next(-1) // -1 = todos los bytes disponibles, sin copia
	if err != nil || len(buf) == 0 {
		return gnet.None
	}

	msgType, derr := protocol.Decode(buf)
	if derr != nil {
		return gnet.None
	}

	if msgType == protocol.TypeStimulus {
		// Responder 1 byte desde el buffer pre-allocated (BR-03).
		if _, werr := c.Write(h.response[:]); werr != nil {
			h.log.Error("write failed", werr)
		}
	} else {
		// Tipo desconocido: ignorar y loggear, no derribar el event loop (BR-01).
		h.log.Info("unknown message type", "type", msgType)
	}
	return gnet.None
}
