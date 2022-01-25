
# Kyoto Benchmarks

Please note, kyoto is not focused on performance (by trying to keep it as performant as possible).
Providing asynchronous extensibility have it's own big overhead.
This benchmarks was made by maintainers, for maintainers.

Benchmarking machine: Macbook Pro 16 (2019)

```
goos: darwin
goarch: amd64
pkg: github.com/kyoto-framework/kyoto/bench
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkEmpty-12                          45788             25301 ns/op
BenchmarkComponents/1-12                   17362             69359 ns/op
BenchmarkComponents/100-12                   370           3174289 ns/op
BenchmarkComponents/1000-12                    8         141959378 ns/op
BenchmarkAction-12                         12346             97740 ns/op
PASS
ok      github.com/kyoto-framework/kyoto/bench  9.624s
```
