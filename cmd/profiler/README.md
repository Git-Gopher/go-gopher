# Profiler

This is setup to profile the main.go application

## How to use

Run the profiler it will generate `cpu.prof` and `mem.prof`
```bash
go run ./cmd/profiler
```

Visualise results with `go tool pprof`
```bash
go tool pprof cpu.prof
# or
go tool pprof mem.prof
```

Example
```
Type: cpu
Time: Aug 2, 2022 at 7:09pm (NZST)
Duration: 75.98s, Total samples = 96.65s (127.21%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) 
```

Type `web` for svg on browser
```
Type: cpu
Time: Aug 2, 2022 at 7:09pm (NZST)
Duration: 75.98s, Total samples = 96.65s (127.21%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) web
```


