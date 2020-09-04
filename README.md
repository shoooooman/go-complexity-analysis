# go-complexity-analysis
go-complexity-analysis calculates the Cyclomatic complexities, the Halstead complexities and the Maintainability indices of golang functions.

# Install
```sh
$ go get github.com/shoooooman/go-complexity-analysis/cmd/complexity
```

# Usage
```sh
$ go vet -vettool=$(which complexity) [flags] [packages]
```

## Flags
`--cycloover`: show functions with the Cyclomatic complexity > N (default: 10)

`--maintunder`: show functions with the Maintainability index < N (default: 20)

## Output
```
<complexity kind> <value> <pkgname> <funcname>
```

## Examples
```go
$ go vet -vettool=$(which complexity) --cycloover 10 .
$ go vet -vettool=$(which complexity) --maintunder 20 main.go
$ go vet -vettool=$(which complexity) --cycloover 5 --maintunder 30 ./src
```
# Metrics
## Cyclomatic Complexity
The Cyclomatic complexity indicates the complexity of a program.
This program calculates the complexities of each function by counting idependent paths with the following rules.
```
Initial value: 1
+1: if, for, case, ||, &&
```

## Halstead Metrics

Reference: https://www.verifysoft.com/en_halstead_metrics.html

### Operands (WIP)
- Identifiers
- Constant
- Typename

### Operators (WIP)
- Operators
- Keywords

## Maintainability Index
The Maintainability index represents maintainability of a program.
The value is calculated with the Cyclomatic complexity and the Halstead volume by using the following formula.
```
Maintainability Index = 171 - 5.2 * ln(Halstead Volume) - 0.23 * (Cyclomatic Complexity) - 16.2 * ln(Lines of Code)
```

This program shows normalized values instead of the original ones [introduced by Microsoft](https://docs.microsoft.com/en-us/archive/blogs/codeanalysis/maintainability-index-range-and-meaning).
```
Normalized Maintainability Index = MAX(0,(171 - 5.2 * ln(Halstead Volume) - 0.23 * (Cyclomatic Complexity) - 16.2 * ln(Lines of Code))*100 / 171)
```

The thresholds are as follows:
```
0-9 = Red
10-19 = Yellow
20-100 = Green
```


# WIP
- [x] Implement the Halstead complexities
- [ ] Connect with Github actions
    - [ ] gets diffs of pull requests
    - [ ] searches functions having the diffs
    - [ ] calculates the difficulties of the pull requests
    - [ ] showes the difficulties
