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