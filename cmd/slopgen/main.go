package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openai/slopgen/internal/check"
	"github.com/openai/slopgen/internal/gen"
)

type result struct {
	index  int
	source []byte
	err    error
}

func main() {
	var seed uint64
	var count, workers, functions, statements, depth, checkEvery int
	var out, pkg string
	flag.Uint64Var(&seed, "seed", 1, "root generation seed")
	flag.IntVar(&count, "count", 1, "number of files to generate")
	flag.IntVar(&workers, "workers", runtime.GOMAXPROCS(0), "parallel workers")
	flag.IntVar(&functions, "functions", 8, "functions per file")
	flag.IntVar(&statements, "statements", 40, "statements per function")
	flag.IntVar(&depth, "max-depth", 4, "maximum expression recursion depth")
	flag.IntVar(&checkEvery, "check-every", 0, "type-check every Nth file (0 disables)")
	flag.StringVar(&out, "out", "-", "output directory, or - for one file on stdout")
	flag.StringVar(&pkg, "package", "generated", "generated package name")
	flag.Parse()
	if count < 1 || workers < 1 || functions < 1 || statements < 1 || depth < 1 {
		fatal("numeric options must be positive")
	}
	if out == "-" && count != 1 {
		fatal("-out=- requires -count=1")
	}
	if out != "-" {
		if err := os.MkdirAll(out, 0o755); err != nil {
			fatal(err.Error())
		}
	}
	cfg := gen.Config{Package: pkg, Functions: functions, Statements: statements, MaxDepth: depth}
	start := time.Now()
	jobs := make(chan int)
	results := make(chan result, workers)
	var total atomic.Uint64
	var wg sync.WaitGroup
	for range min(workers, count) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				src := gen.Generate(seed+uint64(i), cfg)
				var err error
				if checkEvery > 0 && i%checkEvery == 0 {
					err = check.Source(fmt.Sprintf("generated_%06d.go", i), src)
				}
				total.Add(uint64(len(src)))
				results <- result{i, src, err}
			}
		}()
	}
	go func() {
		for i := 0; i < count; i++ {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()
	for item := range results {
		if item.err != nil {
			fatal(fmt.Sprintf("seed %d: %v", seed+uint64(item.index), item.err))
		}
		if out == "-" {
			_, item.err = os.Stdout.Write(item.source)
		} else {
			item.err = os.WriteFile(filepath.Join(out, "generated_"+pad(item.index)+".go"), item.source, 0o644)
		}
		if item.err != nil {
			fatal(item.err.Error())
		}
	}
	elapsed := time.Since(start)
	n := total.Load()
	fmt.Fprintf(os.Stderr, "generated %d files, %.2f MiB in %s (%.2f MiB/s)\n", count, float64(n)/(1<<20), elapsed.Round(time.Millisecond), float64(n)/(1<<20)/elapsed.Seconds())
}
func pad(i int) string {
	s := strconv.Itoa(i)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}
func fatal(message string) { fmt.Fprintln(os.Stderr, "slopgen:", message); os.Exit(2) }
