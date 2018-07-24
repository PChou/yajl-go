# yajl-go

go binding for [yajl](http://lloyd.github.com/yajl), which is a pure ANSI C implementation of json generator and parser

in this binding, `JsonObject` is introduced, and that can be compared:

```go
j5, _ := ParseJson(`{"namelist":[{"name":"jack"},{"name":"rose"}]}`)
j6, _ := ParseJson(`{"namelist":[{"name":"rose"},{"name":"jack"}]}`)
assert.True(t, j5.Compare(j6, DEFAULT))
```

`cgo` must be enabled