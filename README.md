### Book by Ricardo Gerardi


### Lessons Learned

- Use golden file(s) when writing tests.
- `cmd.Run()` is good for interacting with outside processes.
- `flag` is sufficient for small CLI apps, `cobra` is also good.
- Table Driven Tests provide good coverage for most part.
- Don't use `ioutil` because its depreciated. Use `io` or `os` pkg.
- `testdata` directory is ignored when compiling go code.
- use packaged `go tool pprof` for profiling go apps and `go tool trace` for tracing.
- `syscall` is low level library for interacting with outside processes, `os/exec` is high level. 