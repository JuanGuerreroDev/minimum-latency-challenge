package logger

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/JuanGuerreroDev/minimum-latency-challenge/internal/stats"
)

// TestRecordRoundtrip verifica la propiedad de roundtrip (PBT-02): los datos
// registrados con Record aparecen en el archivo escrito por FlushToFile.
func TestRecordRoundtrip(t *testing.T) {
	bl := NewLatencyRecorder(3)

	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	samples := []time.Duration{
		120 * time.Microsecond,
		95 * time.Microsecond,
		310 * time.Microsecond,
	}
	for i, lat := range samples {
		send := base.Add(time.Duration(i) * time.Second)
		bl.Record(send, send.Add(lat), lat)
	}

	if bl.Count() != len(samples) {
		t.Fatalf("Count = %d, se esperaba %d", bl.Count(), len(samples))
	}

	durations := make([]time.Duration, len(samples))
	copy(durations, samples)
	s := stats.Calculate(durations)

	out := filepath.Join(t.TempDir(), "benchmark.log")
	if err := bl.FlushToFile(out, s); err != nil {
		t.Fatalf("FlushToFile falló: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("no se pudo leer el archivo: %v", err)
	}
	content := string(data)

	// Cada latencia registrada debe aparecer en el archivo.
	for _, lat := range samples {
		if !strings.Contains(content, lat.String()) {
			t.Errorf("el archivo no contiene la latencia %v", lat)
		}
	}

	// Verificar que se escribieron exactamente len(samples) líneas de datos.
	dataLines := 0
	sc := bufio.NewScanner(strings.NewReader(content))
	for sc.Scan() {
		line := sc.Text()
		if line != "" && !strings.HasPrefix(line, "#") {
			dataLines++
		}
	}
	if dataLines != len(samples) {
		t.Errorf("líneas de datos = %d, se esperaba %d", dataLines, len(samples))
	}
}

// TestRecordExceedsCapacity verifica que Record crece más allá de la capacidad
// inicial sin perder datos.
func TestRecordExceedsCapacity(t *testing.T) {
	bl := NewLatencyRecorder(1)
	now := time.Now()
	for i := 0; i < 5; i++ {
		bl.Record(now, now, time.Microsecond)
	}
	if bl.Count() != 5 {
		t.Fatalf("Count = %d, se esperaba 5 tras exceder capacidad", bl.Count())
	}
}
