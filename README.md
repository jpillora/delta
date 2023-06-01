# delta

Delta copy from one struct to another and produce a JSON Patch (RFC6902) of the changes

[![GoDev](https://img.shields.io/static/v1?label=godoc&message=reference&color=00add8)](https://pkg.go.dev/github.com/jpillora/deltacopy)
[![CI](https://github.com/jpillora/deltacopy/workflows/CI/badge.svg)](https://github.com/jpillora/deltacopy/actions?workflow=CI)

### Usage

```go
dst := &MyType{...}
src := &MyType{...}

patch, err := delta.CopyPatch(dst, src)
if err != nil {
	return err
}
// dst is now JSON-equivalent to src
// patch is [...RFC6902 operations...]
```

### TODO

* Optimise slice diff