> # Working with errors
>
> The research result.

## Benchmarks

- `make bench`

## Error wrappers

### Candidates

- [github.com/juju/errors](https://github.com/juju/errors/)
- [github.com/pkg/errors](https://github.com/pkg/errors/)

### Summary

| Criterion        | juju/errors | pkg/errors   |
|:-----------------|:-----------:|:------------:|
| Liveliness       | ✅          | ✅           |
| No external deps | ✅          | ✅           |
| Simplicity       | ❎          | ✅           |
| Replace          |             |              |
| - `errors.New`   | ✅          | ✅           |
| - `fmt.Errorf`   | ✅          | ✅           |
| Only context     | ❎          | ✅           |
| Only stack trace | ❎          | ✅           |
| Combination      | ✅          | ✅           |
| Underlying error | ✅          | ✅           |
| License          | LGPLv3      | BSD-2-Clause |
| Score            | ⭐️⭐️⭐️⭐️    | ⭐️⭐️⭐️⭐️     |

### Recommendation

[github.com/juju/errors](https://github.com/juju/errors/) is more complex and rich,
but [github.com/pkg/errors](https://github.com/pkg/errors/) is my preferred choice.

## Panic handlers

### Candidates

- [github.com/oxequa/grace](https://github.com/oxequa/grace/)

### Summary

| Criterion        | oxequa/grace |
|:-----------------|:------------:|
| Liveliness       | ✅           |
| No external deps | ✅           |
| Simplicity       | ✅           |
| Underlying error | ❎           |
| License          | GPLv3        |
| Score            | ⭐️⭐️⭐️ ️     |

### Recommendation

[github.com/oxequa/grace](https://github.com/oxequa/grace/) has some disadvantages,
for example, you cannot obtain the original error and by default the original error
message is concatenated with stack trace without any new lines or spaces.
