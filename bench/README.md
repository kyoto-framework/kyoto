
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

## 10 Apr 2022

```
goos: darwin
goarch: amd64
pkg: github.com/kyoto-framework/kyoto/bench
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkEmpty-12                                          43446             26636 ns/op
BenchmarkComponents/1-12                                   15571             75940 ns/op
BenchmarkComponents/100-12                                   315           3502386 ns/op
BenchmarkComponents/1000-12                                    7         151133768 ns/op
BenchmarkComponentsDynamicRender/1-12                       9656            122647 ns/op
BenchmarkComponentsDynamicRender/100-12                       49          22011405 ns/op
BenchmarkComponentsDynamicRender/1000-12                       1        1811953962 ns/op
BenchmarkComponentsWriter/1-12                             18330             65226 ns/op
BenchmarkComponentsWriter/100-12                             862           1328411 ns/op
BenchmarkComponentsWriter/1000-12                             25          43925223 ns/op
BenchmarkAction-12                                          9584            125116 ns/op
PASS
ok      github.com/kyoto-framework/kyoto/bench  18.495s
```

## 1 May 2022

```
goos: darwin
goarch: amd64
pkg: github.com/kyoto-framework/kyoto/bench
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkEmptyReference-12                                 85660             13356 ns/op
BenchmarkEmpty-12                                          45922             26344 ns/op
BenchmarkComponents/1-12                                   18457             67487 ns/op
BenchmarkComponents/100-12                                   364           3245354 ns/op
BenchmarkComponents/1000-12                                    9         122797136 ns/op
BenchmarkComponentsDynamicRender/1-12                       9266            112225 ns/op
BenchmarkComponentsDynamicRender/100-12                       50          22879164 ns/op
BenchmarkComponentsDynamicRender/1000-12                       1        1779550938 ns/op
BenchmarkComponentsWriter/1-12                             22575             52566 ns/op
BenchmarkComponentsWriter/100-12                            1262            901406 ns/op
BenchmarkComponentsWriter/1000-12                             73          16137558 ns/op
BenchmarkAction-12                                          9339            118158 ns/op
PASS
ok      github.com/kyoto-framework/kyoto/bench  21.057s
```
