package stats

import (
	"testing"
	"time"

	"pgregory.net/rapid"
)

// TestPropertyOrderingInvariant verifica la invariante de ordenamiento (PBT-03):
// Min ≤ Median ≤ Max para cualquier conjunto no vacío de mediciones.
func TestPropertyOrderingInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		raw := rapid.SliceOfN(rapid.Int64Range(0, 1_000_000_000), 1, 10000).Draw(t, "durations")
		durations := make([]time.Duration, len(raw))
		for i, v := range raw {
			durations[i] = time.Duration(v)
		}

		s := Calculate(durations)

		if !(s.Min <= s.Median && s.Median <= s.Max) {
			t.Fatalf("invariante rota: Min(%v) ≤ Median(%v) ≤ Max(%v)", s.Min, s.Median, s.Max)
		}
	})
}

// TestPropertyPercentileInvariant verifica la invariante de percentiles (PBT-03):
// P50 ≤ P95 ≤ P99 para cualquier conjunto no vacío de mediciones.
func TestPropertyPercentileInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		raw := rapid.SliceOfN(rapid.Int64Range(0, 1_000_000_000), 1, 10000).Draw(t, "durations")
		durations := make([]time.Duration, len(raw))
		for i, v := range raw {
			durations[i] = time.Duration(v)
		}

		s := Calculate(durations)

		if !(s.P50 <= s.P95 && s.P95 <= s.P99) {
			t.Fatalf("invariante rota: P50(%v) ≤ P95(%v) ≤ P99(%v)", s.P50, s.P95, s.P99)
		}
	})
}

// TestPropertyCountInvariant verifica la invariante de conteo (PBT-03):
// Count == len(input).
func TestPropertyCountInvariant(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		raw := rapid.SliceOfN(rapid.Int64Range(0, 1_000_000_000), 1, 10000).Draw(t, "durations")
		durations := make([]time.Duration, len(raw))
		for i, v := range raw {
			durations[i] = time.Duration(v)
		}

		s := Calculate(durations)

		if s.Count != len(durations) {
			t.Fatalf("Count = %d, se esperaba %d", s.Count, len(durations))
		}
	})
}

// TestPropertyMedianEqualsP50 verifica que Median == P50 siempre.
func TestPropertyMedianEqualsP50(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		raw := rapid.SliceOfN(rapid.Int64Range(0, 1_000_000_000), 1, 10000).Draw(t, "durations")
		durations := make([]time.Duration, len(raw))
		for i, v := range raw {
			durations[i] = time.Duration(v)
		}

		s := Calculate(durations)
		if s.Median != s.P50 {
			t.Fatalf("Median(%v) != P50(%v)", s.Median, s.P50)
		}
	})
}

// TestPropertyDoesNotMutateInput verifica que Calculate no muta el slice del caller.
func TestPropertyDoesNotMutateInput(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		raw := rapid.SliceOfN(rapid.Int64Range(0, 1_000_000_000), 2, 1000).Draw(t, "durations")
		durations := make([]time.Duration, len(raw))
		original := make([]time.Duration, len(raw))
		for i, v := range raw {
			durations[i] = time.Duration(v)
			original[i] = time.Duration(v)
		}

		_ = Calculate(durations)

		for i := range durations {
			if durations[i] != original[i] {
				t.Fatalf("Calculate mutó el input en índice %d", i)
			}
		}
	})
}

// TestCalculateEmpty verifica el caso de entrada vacía.
func TestCalculateEmpty(t *testing.T) {
	s := Calculate(nil)
	if s.Count != 0 || s.Min != 0 || s.Max != 0 || s.P99 != 0 {
		t.Fatalf("Calculate(nil) debe retornar Stats en cero, obtuvo %+v", s)
	}
}

// TestCalculateKnownValues verifica el cálculo con valores conocidos.
func TestCalculateKnownValues(t *testing.T) {
	durations := []time.Duration{
		5 * time.Microsecond,
		1 * time.Microsecond,
		3 * time.Microsecond,
		2 * time.Microsecond,
		4 * time.Microsecond,
	}
	s := Calculate(durations)

	if s.Min != 1*time.Microsecond {
		t.Errorf("Min = %v, se esperaba 1µs", s.Min)
	}
	if s.Max != 5*time.Microsecond {
		t.Errorf("Max = %v, se esperaba 5µs", s.Max)
	}
	if s.Count != 5 {
		t.Errorf("Count = %d, se esperaba 5", s.Count)
	}
	if s.Avg != 3*time.Microsecond {
		t.Errorf("Avg = %v, se esperaba 3µs", s.Avg)
	}
}
