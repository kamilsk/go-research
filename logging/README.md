> # Logging
>
> The research result.

## Benchmarks

- [Benchmarking logging libraries for Go](https://github.com/imkira/go-loggers-bench)
- [Zap benchmarking suite](https://github.com/uber-go/zap#performance)
- `make bench`

## Structured loggers

### Candidates

- [github.com/sirupsen/logrus][logrus]
- [go.uber.org/zap][zap]
- [github.com/rs/zerolog][zerolog]

### Summary

| Criteria         | [github.com/sirupsen/logrus][logrus] | [go.uber.org/zap][zap] | [github.com/rs/zerolog][zerolog] |
|:-----------------|:---------:|:--------:|:-------:|
| Liveliness       | ✅        | ✅       | ✅      |
| No external deps | ❎        | ❎       | ✅      |
| Simplicity       | ✅        | ❎       | ❎      |
| Substitution     |           |          |         |
| - `log.Logger`   | ✅        | ❎       | ❎      |
| - `io.Writer`    | ✅        | ❎       | ✅      |
| Documentation    |           |          |         |
| - benchmarks     | ✅        | ✅       | ✅      |
| - examples       | ✅        | ❎       | ✅      |
| Features         |           |          |         |
| - formatter      | ✅        | ✅       | ✅      |
| - levels         | ✅        | ✅       | ✅      |
| - nesting        | ✅        | ✅       | ✅      |
| - structured     | ✅        | ✅       | ✅      |
| - timestamp      | ✅        | ✅       | ✅      |
| License          | MIT       | MIT      | MIT     |
| Score            | ⭐⭐⭐⭐  | ⭐⭐⭐⭐  | ⭐⭐⭐  |

### Recommendation

[go.uber.org/zap][zap] has good performance and flexible configuration.

[logrus]:  https://github.com/sirupsen/logrus
[zap]:     https://github.com/uber-go/zap
[zerolog]: https://github.com/rs/zerolog
