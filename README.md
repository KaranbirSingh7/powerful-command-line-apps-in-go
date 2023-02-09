### Book by Ricardo Gerardi


### Lessons Learned

- Use golden file(s) when writing tests.
- `cmd.Run()` is good for interacting with outside processes.
- `flag` is sufficient for small CLI apps, `cobra` is also good.
- Table Driven Tests provide good coverage for most part.
- Don't use `ioutil` because its depreciated. Use `io` or `os` pkg.
- `testdata` directory is ignored when compiling go code.
- when testing use `t.Helper()` for helper functions, `t.Main()` for pre-reqs and post cleanup tasks. also use `t.Parallel()` where possible to speed up test runs. 
- use packaged `go tool pprof` for profiling go apps and `go tool trace` for tracing.
- `syscall` is low level library for interacting with outside processes, `os/exec` is high level. 
- whenever running external commands/processes where long timeout can happen, always use `context`
- use constructors for creating structs, example: `func newSqlClient(opts Options)`. This is where you will apply any struct defaults.


### Lessons Learned - Code 

- Use wrapper struct around when executing external processes
	```go
	# file: step.go
	func execute(exe string, proj string, []string args) (string, error){
		cmd := exec.Command(exe, args...)
		cmd.Dir = proj

		// to capture output from external process
		var out bytes.Buffer
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to run command: %w", err)
		}
		return out, nil
	}
	```

