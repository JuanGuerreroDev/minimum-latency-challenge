// Package protocol define el protocolo binario ultra-minimal del sistema de
// mínima latencia. Un mensaje es exactamente 1 byte (solo el tipo, sin payload),
// lo que elimina todo overhead de serialización (NFR-PERF-03, BR-02).
//
// Las funciones Encode y Decode son puras y zero-allocation: operan sobre
// buffers provistos por el caller y nunca asignan en el heap (PAT-PERF-03).
package protocol

import "errors"

// Tipos de mensaje del protocolo (BR-01).
const (
	// TypeStimulus es el mensaje Cliente → Servidor (0x01).
	TypeStimulus byte = 0x01
	// TypeResponse es el mensaje Servidor → Cliente (0x02).
	TypeResponse byte = 0x02

	// MessageSize es el tamaño en bytes de todo mensaje válido en el wire.
	// El protocolo es ultra-minimal: 1 byte de tipo, sin payload (NFR-PERF-03).
	MessageSize int = 1
)

// ErrEmptyMessage se retorna cuando se intenta decodificar un buffer vacío.
var ErrEmptyMessage = errors.New("protocol: empty message")

// Encode escribe msgType en la primera posición de buf y retorna el número de
// bytes escritos (siempre MessageSize). Es zero-allocation: usa el buffer
// pre-allocated del caller y no asigna memoria (PAT-PERF-03, PAT-PERF-04).
//
// El caller debe garantizar len(buf) >= MessageSize.
func Encode(msgType byte, buf []byte) int {
	buf[0] = msgType
	return MessageSize
}

// Decode lee el tipo de mensaje desde data. Retorna ErrEmptyMessage si data
// está vacío. Es zero-allocation: no copia ni asigna (PAT-PERF-03).
//
// Solo se interpreta el primer byte; los bytes restantes se ignoran (BR-02).
func Decode(data []byte) (msgType byte, err error) {
	if len(data) < MessageSize {
		return 0, ErrEmptyMessage
	}
	return data[0], nil
}

// IsValidType reporta si t es un tipo de mensaje conocido del protocolo (BR-01).
func IsValidType(t byte) bool {
	return t == TypeStimulus || t == TypeResponse
}
