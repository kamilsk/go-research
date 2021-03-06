package strings

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func buildByBuilder(id string) string {
	b, base := strings.Builder{}, "https://guard.octolab.net/api/v1/license/"
	b.Grow(len(base) + len(id))
	_, _ = b.WriteString(base)
	_, _ = b.WriteString(id)
	return b.String()
}

func buildByBuffer(id string) string {
	b, base := bytes.Buffer{}, "https://guard.octolab.net/api/v1/license/"
	b.Grow(len(base) + len(id))
	_, _ = b.WriteString(base)
	_, _ = b.WriteString(id)
	return b.String()
}

func buildByConcat(id string) string {
	return "https://guard.octolab.net/api/v1/license/" + id
}

func buildByMask(id string) string {
	base := []byte("https://guard.octolab.net/api/v1/license/00000000-0000-0000-0000-000000000000")
	window := base[len(base)-len(id):]
	_ = copy(window, []byte(id))
	return string(base)
}

func buildByPrint(id string) string {
	return fmt.Sprintf("https://guard.octolab.net/api/v1/license/%s", id)
}

func buildByUnsafe(id string) string {
	base := "https://guard.octolab.net/api/v1/license/00000000-0000-0000-0000-000000000000"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&base))
	data := *(*[77]byte)(unsafe.Pointer(hdr.Data))
	hdr.Data = uintptr(unsafe.Pointer(&data))
	window := data[len(data)-len(id):]
	_ = copy(window, id)
	return base
}

// BenchmarkBuild/build_by_builder-12         	 3000000	       421 ns/op	     800 B/op	      10 allocs/op
// BenchmarkBuild/build_by_buffer-12          	 1000000	      1876 ns/op	    2720 B/op	      30 allocs/op
// BenchmarkBuild/build_by_concat-12          	 3000000	       466 ns/op	     800 B/op	      10 allocs/op
// BenchmarkBuild/build_by_mask-12            	 1000000	      1031 ns/op	    2080 B/op	      30 allocs/op
// BenchmarkBuild/build_by_print-12           	 1000000	      1386 ns/op	     960 B/op	      20 allocs/op
// BenchmarkBuild/build_by_unsafe-12          	 5000000	       337 ns/op	     800 B/op	      10 allocs/op
// compare with https://gist.github.com/kamilsk/af63aa5bb6178d4e4aeef091bdf32696

func BenchmarkBuild(b *testing.B) {
	ids := []string{
		"072f5a3c-d612-4b91-af13-937e36e2aa93",
		"3303ac29-8982-4220-a405-c5694b3f0baf",
		"43fee265-a980-4bce-bb8b-8e780c67e047",
		"9b8d8fdb-dd04-4632-beff-261b132231a0",
		"8249a5cf-7bfb-4107-9bc0-f1a17dfa90e5",
		"6b3b2eb5-478f-4746-97f8-3a2f2c5b8aef",
		"d8de660b-739b-4bed-a961-df39822bbd3b",
		"1bdb82bc-0077-404d-bf2a-dda69ded1ce4",
		"bff662fd-95a6-4f60-92d9-a052bee6779c",
		"838dc73f-983c-4cac-8690-252a4ea4bbaa",
	}

	benchmarks := []struct {
		name      string
		algorithm func(string) string
	}{
		{"build by builder", buildByBuilder},
		{"build by buffer", buildByBuffer},
		{"build by concat", buildByConcat},
		{"build by mask", buildByMask},
		{"build by print", buildByPrint},
		{"build by unsafe", buildByUnsafe},
	}
	for _, bm := range benchmarks {
		tc := bm
		b.Run(bm.name, func(b *testing.B) {
			var last string
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, id := range ids {
					last = tc.algorithm(id)
				}
			}
			if last != "https://guard.octolab.net/api/v1/license/838dc73f-983c-4cac-8690-252a4ea4bbaa" {
				b.Fail()
			}
		})
	}
}
