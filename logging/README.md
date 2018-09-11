> # Logging
>
> The research result.

## Benchmarks

- [Benchmarking logging libraries for Go](https://github.com/imkira/go-loggers-bench)
- [Zap benchmarking suite](https://github.com/uber-go/zap#performance)
- `make bench`

## Structured loggers

### Candidates

- [github.com/sirupsen/logrus](https://github.com/Sirupsen/logrus)
- [go.uber.org/zap](https://github.com/uber-go/zap)
- [github.com/rs/zerolog](https://github.com/rs/zerolog)

### Summary

| Criteria         | github.com/sirupsen/logrus | go.uber.org/zap | github.com/rs/zerolog |
|:-----------------|:--------------------------:|:---------------:|:---------------------:|
| Liveliness       | ✅                         | ✅              | ✅                    |
| No external deps | ❎                         | ❎              | ✅                    |
| Simplicity       | ✅                         | ❎              | ❎                    |
| Compatible       |                            |                 |                       |
| - `log.Logger`   | ✅                         | ❎              | ❎                    |
| - `io.Writer`    | ✅                         | ❎              | ✅                    |
| Documentation    |                            |                 |                       |
| - benchmarks     | ✅                         | ✅              | ✅                    |
| - examples       | ✅                         | ❎              | ✅                    |
| Features         |                            |                 |                       |
| - formatter      | ✅                         | ✅              | ✅                    |
| - levels         | ✅                         | ✅              | ✅                    |
| - nesting        | ✅                         | ✅              | ✅                    |
| - structured     | ✅                         | ✅              | ✅                    |
| - timestamp      | ✅                         | ✅              | ✅                    |
| License          | MIT                        | MIT             | MIT                   |
| Score            | ⭐⭐⭐⭐                   | ⭐⭐⭐⭐         | ⭐⭐⭐                |

### Recommendation

- [ ] TODO need more criteria
- [ ] TODO complete benchmarks for usual cases

### Next iteration

- [github.com/inconshreveable/log15](https://github.com/inconshreveable/log15)
- [github.com/apex/log](https://github.com/apex/log)
- [github.com/francoispqt/onelog](https://github.com/francoispqt/onelog)
