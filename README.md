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
## Halstead Metrics

Calculation of each Halstead metrics can be found [here](!https://www.verifysoft.com/en_halstead_metrics.html).

### Rules
1. Comments are not considered in Halstead Metrics
2. Operands and Operators are divided as follows:
#### Operands 
- [Identifiers](!https://golang.org/ref/spec#Identifiers)
- [Constants](!https://golang.org/ref/spec#Constants)
- [Variables](!https://golang.org/ref/spec#Variables)

#### Operators
- [Operators](!https://golang.org/ref/spec#Operators_and_punctuation)  
    - Parenthesis, such as "()", is counted as one operator
- [Keywords](!https://golang.org/ref/spec#Keywords)


# WIP
- [x] Implement the Halstead complexities
- [ ] Connect with Github actions
    - [ ] gets diffs of pull requests
    - [ ] searches functions having the diffs
    - [ ] calculates the difficulties of the pull requests
    - [ ] showes the difficulties
