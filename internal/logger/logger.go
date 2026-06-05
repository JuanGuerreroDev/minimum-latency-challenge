// Package logger provee logging estructurado (JSON) para el sistema y un
// BenchmarkLogger buffered para la trazabilidad de mediciones de latencia.
//
// El logger de sistema se basa en log/slog (stdlib, Go 1.21+) y emite JSON con
// timestamp ISO 8601, level y message (NFR-LOG-01, PAT-OBS-01, SECURITY-03).
// Solo se usa en cold paths; nunca en el éxito del hot path del servidor.
package logger

import (
	"io"
	"log/slog"
)

// Logger envuelve un *slog.Logger configurado para salida JSON estructurada.
type Logger struct {
	*slog.Logger
}

// New crea un Logger estructurado que escribe JSON a output.
// El nivel mínimo es Info. El timestamp se emite en formato ISO 8601 (RFC 3339).
func New(output io.Writer) *Logger {
	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return &Logger{Logger: slog.New(handler)}
}

// Error registra un mensaje de nivel ERROR adjuntando el error como campo "err".
// Acepta pares clave-valor adicionales (estilo slog).
func (l *Logger) Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("err", err))
	}
	l.Logger.Error(msg, args...)
}
