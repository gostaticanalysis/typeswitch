# typeswitch

[![godoc.org][godoc-badge]][godoc]

`typeswitch` finds a type which implement an interfaces which are used in type-switch but the type does not appear in any cases of the type-switch.

```go
package main

type I interface{ F() }
type A struct{I} // implements I
type B struct{I} // implements I

func main() {
	var i I = A{}
	switch i.(type) {
	case A:
	}
}
```

```sh
$go vet -vettool=`which typeswitch` main.go
./main.go:9:2: type B does not appear in any cases
```

<!-- links -->
[godoc]: https://godoc.org/github.com/gostaticanalysis/typeswitch
[godoc-badge]: https://img.shields.io/badge/godoc-reference-4F73B3.svg?style=flat-square&label=%20godoc.org
