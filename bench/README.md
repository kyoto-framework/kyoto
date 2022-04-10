
# Kyoto Benchmarks

Please note, kyoto is not focused on performance (by trying to keep it as performant as possible).
Providing asynchronous extensibility have it's own big overhead.
This benchmarks was made by maintainers, for maintainers.

## 25 Jan 2022

Benchmarking machine: Macbook Pro 16 (2019)

```
goos: darwin
goarch: amd64
pkg: github.com/kyoto-framework/kyoto/bench
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkEmpty-12                          44184             25820 ns/op
BenchmarkComponents/1-12                   16381             72998 ns/op
BenchmarkComponents/100-12                   345           3617490 ns/op
BenchmarkComponents/1000-12                    6         181818836 ns/op
BenchmarkAction-12                          7335            155638 ns/op
PASS
ok      github.com/kyoto-framework/kyoto/bench  8.390s
```

## 9 Apr 2022

```
goos: darwin
goarch: amd64
pkg: github.com/kyoto-framework/kyoto/bench
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkEmpty-12                          46172             25432 ns/op
BenchmarkComponents/1-12                   16909             70940 ns/op
BenchmarkComponents/100-12                   339           3359374 ns/op
BenchmarkComponents/1000-12                    7         148261588 ns/op
BenchmarkAction-12                          9643            125656 ns/op
PASS
ok      github.com/kyoto-framework/kyoto/bench  9.958s
```
