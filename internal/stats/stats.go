// Package stats calcula estadísticas de latencia a partir de las mediciones del
// benchmark: mínimo, máximo, promedio, mediana y percentiles p50/p95/p99.
//
// Calculate es una función pura: ordena una copia local de las duraciones y
// computa las métricas. Las estadísticas se calculan SOLO sobre iteraciones
// exitosas (BR-05).
package stats

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Stats contiene las estadísticas de latencia calculadas sobre un conjunto de
// mediciones exitosas.
//
// Invariantes (PAT-TEST-01):
//   - Min ≤ Median ≤ Max
//   - P50 ≤ P95 ≤ P99
//   - P50 == Median
//   - Avg == Total / Count
type Stats struct {
	Count  int           // Número de mediciones exitosas
	Min    time.Duration // Latencia mínima
	Max    time.Duration // Latencia máxima
	Avg    time.Duration // Latencia promedio
	Median time.Duration // Mediana (== P50)
	P50    time.Duration // Percentil 50
	P95    time.Duration // Percentil 95
	P99    time.Duration // Percentil 99
	Total  time.Duration // Suma total de todas las latencias
}

// percentileIndex retorna el índice (0-based) correspondiente a un percentil
// dado sobre un slice ordenado de longitud n. Usa el método nearest-rank
// acotado para no salir de rango.
func percentileIndex(percentile, n int) int {
	if n == 0 {
		return 0
	}
	idx := (percentile * n) / 100
	if idx >= n {
		idx = n - 1
	}
	return idx
}

// Calculate computa las estadísticas a partir de las duraciones medidas.
// Si durations está vacío, retorna un Stats con todos los campos en cero.
//
// No muta el slice del caller: ordena una copia local.
func Calculate(durations []time.Duration) *Stats {
	n := len(durations)
	if n == 0 {
		return &Stats{}
	}

	sorted := make([]time.Duration, n)
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	var total time.Duration
	for _, d := range sorted {
		total += d
	}

	return &Stats{
		Count:  n,
		Min:    sorted[0],
		Max:    sorted[n-1],
		Avg:    total / time.Duration(n),
		Median: sorted[percentileIndex(50, n)],
		P50:    sorted[percentileIndex(50, n)],
		P95:    sorted[percentileIndex(95, n)],
		P99:    sorted[percentileIndex(99, n)],
		Total:  total,
	}
}

// String retorna una representación compacta de las estadísticas en una línea.
func (s *Stats) String() string {
	return fmt.Sprintf(
		"count=%d min=%v max=%v avg=%v p50=%v p95=%v p99=%v",
		s.Count, s.Min, s.Max, s.Avg, s.P50, s.P95, s.P99,
	)
}

// Report retorna un reporte tabular detallado para imprimir en stdout.
// Incluye la comparación contra el objetivo de latencia p99 < 1ms.
func (s *Stats) Report() string {
	var b strings.Builder
	b.WriteString("================ Latency Benchmark Results ================\n")
	fmt.Fprintf(&b, "  Successful iterations : %d\n", s.Count)
	fmt.Fprintf(&b, "  Min                   : %v\n", s.Min)
	fmt.Fprintf(&b, "  Max                   : %v\n", s.Max)
	fmt.Fprintf(&b, "  Avg                   : %v\n", s.Avg)
	fmt.Fprintf(&b, "  Median (p50)          : %v\n", s.Median)
	fmt.Fprintf(&b, "  p95                   : %v\n", s.P95)
	fmt.Fprintf(&b, "  p99                   : %v\n", s.P99)
	b.WriteString("-----------------------------------------------------------\n")

	const target = time.Millisecond
	if s.Count > 0 && s.P99 < target {
		fmt.Fprintf(&b, "  Objetivo p99 < 1ms    : ✅ ALCANZADO (p99 = %v)\n", s.P99)
	} else {
		fmt.Fprintf(&b, "  Objetivo p99 < 1ms    : ❌ NO ALCANZADO (p99 = %v)\n", s.P99)
	}
	b.WriteString("===========================================================\n")
	return b.String()
}
