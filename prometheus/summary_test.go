// Copyright 2014 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"sync"
	"testing"

	dto "github.com/prometheus/client_model/go"
)

func benchmarkSummaryObserve(w int, b *testing.B) {
	b.StopTimer()

	wg := new(sync.WaitGroup)
	wg.Add(w)

	g := new(sync.WaitGroup)
	g.Add(1)

	s := NewSummary(SummaryOpts{})

	for i := 0; i < w; i++ {
		go func() {
			g.Wait()

			for i := 0; i < b.N; i++ {
				s.Observe(float64(i))
			}

			wg.Done()
		}()
	}

	b.StartTimer()
	g.Done()
	wg.Wait()
}

func BenchmarkSummaryObserve1(b *testing.B) {
	benchmarkSummaryObserve(1, b)
}

func BenchmarkSummaryObserve2(b *testing.B) {
	benchmarkSummaryObserve(2, b)
}

func BenchmarkSummaryObserve4(b *testing.B) {
	benchmarkSummaryObserve(4, b)
}

func BenchmarkSummaryObserve8(b *testing.B) {
	benchmarkSummaryObserve(8, b)
}

func benchmarkSummaryWrite(w int, b *testing.B) {
	b.StopTimer()

	wg := new(sync.WaitGroup)
	wg.Add(w)

	g := new(sync.WaitGroup)
	g.Add(1)

	s := NewSummary(SummaryOpts{})

	for i := 0; i < 1000000; i++ {
		s.Observe(float64(i))
	}

	for j := 0; j < w; j++ {
		outs := make([]dto.Metric, b.N)

		go func(o []dto.Metric) {
			g.Wait()

			for i := 0; i < b.N; i++ {
				s.Write(&o[i])
			}

			wg.Done()
		}(outs)
	}

	b.StartTimer()
	g.Done()
	wg.Wait()
}

func BenchmarkSummaryWrite1(b *testing.B) {
	benchmarkSummaryWrite(1, b)
}

func BenchmarkSummaryWrite2(b *testing.B) {
	benchmarkSummaryWrite(2, b)
}

func BenchmarkSummaryWrite4(b *testing.B) {
	benchmarkSummaryWrite(4, b)
}

func BenchmarkSummaryWrite8(b *testing.B) {
	benchmarkSummaryWrite(8, b)
}

func ExampleSummary() {
	temps := NewSummary(SummaryOpts{
		Name: "pond_temperature",
		Help: "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
	})

	temps.Observe(37)
	// - count:   1
	// - sum:    37
	// - median: 37
	// - 90th:   37
	// - 99th:   37
}

func ExampleSummaryVec() {
	temps := NewSummaryVec(
		SummaryOpts{
			Name: "pond_temperature",
			Help: "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
		},
		[]string{"species"},
	)

	temps.WithLabelValues("litoria-caerulea").Observe(37) // Not so stinky.

	temps.WithLabelValues("lithobates-catesbeianus").Observe(40) // Quite stinky!
	// Grab a beer to drown away the pain of the smell before sampling again.
	temps.WithLabelValues("lithobates-catesbeianus").Observe(42)
	// species: litoria-caerulea
	// - count:   1
	// - sum:    37
	// - median: 37
	// - 90th:   37
	// - 99th:   37
	// species: lithobates-catesbeianus
	// - count:   2
	// - sum:    82
	// - median: 41
	// - 90th:   42
	// - 99th:   42
}