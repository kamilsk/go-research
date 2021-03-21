> # Working with errors
>
> The research result.

## Benchmarks

- `make bench`

## Error wrappers

### Candidates

- [github.com/goph/emperror][emperror]
- [github.com/juju/errors][juju/errors]
- [github.com/pkg/errors][pkg/errors]

### Summary

| Criterion        | [emperror][] | [juju/errors][] | [pkg/errors][] |
|:-----------------|:---:|:---:|:---:|
| Liveliness       | ✅ | ✅ | ✅ |
| No external deps | ❎ | ✅ | ✅ |
| Simplicity       | ✅ | ❎ | ✅ |
| Substitution     |
| - `errors.New`   |   | ✅ | ✅ |
| - `fmt.Errorf`   |   | ✅ | ✅ |
| Packing with     |
| - context        |   | ❎ | ❎ |
| - description    |   | ❎ | ✅ |
| - error          |   | ✅ | ✅ |
| - stack trace    |   | ❎ | ✅ |
| Unpacking        |
| - context        |
| - description    |
| - error          |
| - stack trace    |
| Underlying error | ✅ | ✅ | ✅ |
| License          | LGPLv3 | BSD-2-Clause | MIT |
| Score            | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

<sup id="anchor-">[1](#)</sup>
<sup id="">1</sup> [↩](#anchor-)

### Recommendation

[github.com/juju/errors][juju/errors] is more complex and rich, but [github.com/pkg/errors][pkg/errors] is
my preferred choice because of simplicity.

[github.com/goph/emperror][emperror] extends [github.com/pkg/errors][pkg/errors] and provides useful ideas
and integrations, but has not optimal code, `emperror.Context()` as example, and copy/paste from [pkg/errors][].

## Panic handlers

### Candidates

- [github.com/goph/emperror][emperror]
- [github.com/oxequa/grace][grace]

### Summary

| Criterion        | [emperror][] |[grace][] |
|:-----------------|:---:|:---:|
| Liveliness       | ✅ | ✅ |
| No external deps | ❎ | ✅ |
| Simplicity       | ✅ | ✅ |
| Underlying error | ✅ | ❎ |
| License          | MIT | GPLv3 |
| Score            | ⭐⭐⭐⭐ | ⭐⭐⭐ |

### Recommendation

[github.com/oxequa/grace][grace] has some disadvantages, for example, you cannot obtain the original error
and by default the original error message is concatenated with stack trace without any new lines or spaces.

[github.com/goph/emperror][emperror] is closest to my needs, but has an external dependency.

[emperror]:    https://github.com/goph/emperror
[juju/errors]: https://github.com/juju/errors
[pkg/errors]:  https://github.com/pkg/errors
[grace]:       https://github.com/oxequa/grace
