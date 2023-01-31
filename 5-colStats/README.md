## Testing

1. Run all tests
    ```sh
    go test -v ./...
    ```

## Installation

1. Build binary
    ```sh
    go build
    ```

1. Run program
    ```sh
    ./colstat -op avg -col 3 testdata/example.csv
    ```


### Benchmarking

1. Extract sample files to run test with load
    ```sh
    tar -xvzf testdata/colStatsBenchmarkData.tar.gz -C ./testdata/
    ```

1. Run benchmark tests
    ```sh
    go test -bench . -run ^$
    # ^$ is special symbol that ignore regular testCases
    ```

1. Run benchmark `n` times
    ```sh
    go test -bench=. -benchtime=10x -run ^$
    # ^$ is special symbol that ignore regular testCases
    ```


1. Run benchmark `n` times and display total memory consumption
    ```sh
    go test -bench=. -benchtime=10x -run ^$ -benchmem 
    # ^$ is special symbol that ignore regular testCases
    ```


1. Run benchmark with profiler enabled - CPU consumption based
    This will generate (2) files, 1 - your progra compiled binary and 2 - your profiling results that can be analyzed using `go tool pprof`
    ```sh
    go test -bench=. -benchtime=10x -run ^$ -cpuprofile cpu00.pprof
    # ^$ is special symbol that ignore regular testCases

    # analyze using pprof
    go tool pprof cpu00.pprof
    ```


1. Run benchmark with profiler enabled - Memory consumption based
    This will generate (2) files, 1 - your progra compiled binary and 2 - your profiling results that can be analyzed using `go tool pprof`
    ```sh
    go test -bench=. -benchtime=10x -run ^$ -memprofile mem00.pprof
    # ^$ is special symbol that ignore regular testCases

    # analyze using pprof
    go tool pprof mem00.pprof
    ```

