package protocol

import (
	"testing"

	"pgregory.net/rapid"
)

// TestPropertyEncodeDecodeRoundtrip verifica la propiedad de roundtrip
// (PBT-02): Decode(Encode(t)) == t para todo tipo de mensaje (BR, PAT-TEST-01).
func TestPropertyEncodeDecodeRoundtrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		msgType := byte(rapid.IntRange(0, 255).Draw(t, "msgType"))

		buf := make([]byte, MessageSize)
		n := Encode(msgType, buf)

		if n != MessageSize {
			t.Fatalf("Encode escribió %d bytes, se esperaba %d", n, MessageSize)
		}

		decoded, err := Decode(buf[:n])
		if err != nil {
			t.Fatalf("Decode falló sobre output de Encode: %v", err)
		}
		if decoded != msgType {
			t.Fatalf("roundtrip roto: Encode(%#x) → Decode = %#x", msgType, decoded)
		}
	})
}

// TestPropertyEncodeSizeInvariant verifica la invariante de tamaño (PBT-03):
// Encode siempre escribe exactamente 1 byte (NFR-PERF-03).
func TestPropertyEncodeSizeInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		msgType := byte(rapid.IntRange(0, 255).Draw(t, "msgType"))
		buf := make([]byte, 8) // buffer holgado a propósito

		n := Encode(msgType, buf)
		if n != 1 {
			t.Fatalf("invariante de tamaño rota: Encode escribió %d bytes, debe ser 1", n)
		}
	})
}

// TestDecodeEmpty verifica el manejo de mensaje vacío.
func TestDecodeEmpty(t *testing.T) {
	if _, err := Decode(nil); err != ErrEmptyMessage {
		t.Fatalf("Decode(nil) = %v, se esperaba ErrEmptyMessage", err)
	}
	if _, err := Decode([]byte{}); err != ErrEmptyMessage {
		t.Fatalf("Decode([]) = %v, se esperaba ErrEmptyMessage", err)
	}
}

// TestDecodeIgnoresExtraBytes verifica BR-02: solo se interpreta el primer byte.
func TestDecodeIgnoresExtraBytes(t *testing.T) {
	got, err := Decode([]byte{TypeStimulus, 0xFF, 0xAA})
	if err != nil {
		t.Fatalf("Decode falló: %v", err)
	}
	if got != TypeStimulus {
		t.Fatalf("Decode = %#x, se esperaba %#x (solo primer byte)", got, TypeStimulus)
	}
}

// TestIsValidType verifica la validación de tipos conocidos (BR-01).
func TestIsValidType(t *testing.T) {
	cases := map[byte]bool{
		TypeStimulus: true,
		TypeResponse: true,
		0x00:         false,
		0x03:         false,
		0xFF:         false,
	}
	for in, want := range cases {
		if got := IsValidType(in); got != want {
			t.Errorf("IsValidType(%#x) = %v, se esperaba %v", in, got, want)
		}
	}
}

// BenchmarkEncode verifica zero-allocation en el hot path (NFR-PERF-02).
// go test -bench=. -benchmem debe reportar 0 allocs/op.
func BenchmarkEncode(b *testing.B) {
	buf := make([]byte, MessageSize)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Encode(TypeResponse, buf)
	}
}

// BenchmarkDecode verifica zero-allocation en el hot path (NFR-PERF-02).
func BenchmarkDecode(b *testing.B) {
	data := []byte{TypeStimulus}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(data)
	}
}
