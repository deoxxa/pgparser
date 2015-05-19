pgparser
========

[![GoDoc](https://godoc.org/fknsrs.biz/p/pgparser?status.svg)](https://godoc.org/fknsrs.biz/p/pgparser)

Overview
--------

This package can unpack complex types from Postgres into structs, slices, and
primitive values.

Example
-------

Also see `example_test.go` in this directory.

```go
import (
  "fmt"

  "fknsrs.biz/p/pgparser"
)

func ExampleParse() {
  var s = `(a,{x,y,z},5)`

  var v struct{
    A string
    B []string
    C int
  }

  if err := pgparser.Unmarshal(s, &v); err != nil {
    panic(err)
  }

  fmt.Printf("%#v\n", v)

  // Output:
  // struct { A string; B []string; C int }{A:"a", B:[]string{"x", "y", "z"}, C:5}
}
```

License
-------

3-clause BSD. A copy is included with the source.
