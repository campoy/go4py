package lists

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"testing"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func BenchmarkConcrete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterInt(func(x int) bool { return x%2 == 0 },
			MapIntInt(func(x int) int { return x * x },
				Range(0, 1000)))
	}
}

func BenchmarkGeneric(b *testing.B) {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	for i := 0; i < b.N; i++ {
		Filter(func(x int) bool { return x%2 == 0 },
			Map(func(x int) int { return x * x },
				Range(0, 1000)))
	}
}
