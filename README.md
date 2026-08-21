# slopgen

`slopgen` deterministically generates large quantities of type-correct Go source. Its M0 generator walks weighted grammar productions while restricting each expression to the requested scalar type. Compilation is not part of the generation hot path; optional sampled validation uses `go/types` in-process.

## Build and run

```sh
go build ./cmd/slopgen
./slopgen -seed 42 -count 100 -out ./generated -workers 8 -check-every 10
```

The same seed and generation settings produce byte-identical files regardless of worker count. A single file can be written to standard output:

```sh
go run ./cmd/slopgen -seed 42 -out - > generated.go
```

Key controls are `-functions`, `-statements`, and `-max-depth`. Progress metrics are printed to standard error. M0 deliberately emits one import-free package containing scalar expressions, bounded control flow, and functions whose final statement is always a return.
