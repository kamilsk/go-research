package strings

import (
	"sort"
	"testing"
)

func filterByAppend(input []string) []string {
	output := input[:0]
	for _, num := range input {
		if num != "" {
			output = append(output, num)
		}
	}
	return output
}

func filterBySort(input []string) []string {
	var border int
	sort.Strings(input)
	for i, num := range input {
		if num != "" {
			border = i
			break
		}
	}
	return input[border:]
}

func filterBySwap(input []string) []string {
	var zero int
	for i, num := range input {
		if num != "" {
			input[i], input[zero] = input[zero], input[i]
			zero++
		}
	}
	return input[:zero]
}

// BenchmarkFilter/filter_by_append-12         	50000000	        30.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkFilter/filter_by_sort-12           	 2000000	       695 ns/op	      32 B/op	       1 allocs/op
// BenchmarkFilter/filter_by_swap-12           	30000000	        40.5 ns/op	       0 B/op	       0 allocs/op
// compare with https://gist.github.com/kamilsk/a31a7a3fe2a0dc7cbac49ec449f73594

func BenchmarkFilter(b *testing.B) {
	input := []string{
		"072f5a3c-d612-4b91-af13-937e36e2aa93", "",
		"3303ac29-8982-4220-a405-c5694b3f0baf", "",
		"43fee265-a980-4bce-bb8b-8e780c67e047", "",
		"9b8d8fdb-dd04-4632-beff-261b132231a0", "",
		"8249a5cf-7bfb-4107-9bc0-f1a17dfa90e5", "",
		"6b3b2eb5-478f-4746-97f8-3a2f2c5b8aef", "",
		"d8de660b-739b-4bed-a961-df39822bbd3b", "",
		"1bdb82bc-0077-404d-bf2a-dda69ded1ce4", "",
		"bff662fd-95a6-4f60-92d9-a052bee6779c", "",
		"838dc73f-983c-4cac-8690-252a4ea4bbaa", "",
	}

	benchmarks := []struct {
		name      string
		algorithm func([]string) []string
	}{
		{"filter by append", filterByAppend},
		{"filter by sort", filterBySort},
		{"filter by swap", filterBySwap},
	}
	for _, bm := range benchmarks {
		tc := bm
		b.Run(bm.name, func(b *testing.B) {
			var last []string
			copied := make([]string, len(input))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				copy(copied, input)
				last = tc.algorithm(copied)
			}
			if len(last) != 10 {
				b.Fail()
			}
		})
	}
}
