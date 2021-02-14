# A quick tour of Go 1.16 @ 16th Golang Leipzig Meetup

Our 16th Meetup at the 16th of February with a quick tour on Go 1.16 :blush: .
We will also have an outlook on the Go generics proposal.

Go 1.16 is expected to be released this Februrary! :tada: :fireworks:

## From the [preliminary release notes](https://tip.golang.org/doc/go1.16#ioutil)

- :warning: [`io/ioutil` gets deprecated](https://tip.golang.org/doc/go1.16#ioutil)
  - package remains and will continue to work like before but new code should use functions from `io` and `os` instead
  - `ioutil` functions will become wrappers around `io` and `os`, example [`ioutil.ReadFile`](https://tip.golang.org/src/io/ioutil/ioutil.go?s=1167:1213#L26)
- :warning: 32-bit iOS and macOS support is dropped
- :tada: Go adds support for the 64-bit macOS ARM architecture (aka Apple Silicon) with `GOOS=darwin, GOARCH=arm64`
  - ~~crypto will be slower on this architecture since the optimized assembly code is still missing~~ There are arm64 assembly implementations: https://tip.golang.org/src/crypto/aes/cipher_asm.go and https://tip.golang.org/src/crypto/aes/asm_arm64.s

#### [Modules](https://tip.golang.org/doc/go1.16#modules)

- :warning: `GO111MODULE=on` is now the default, i.e. go tools work in  _go modules_ mode by default
- `go build` and `go test` do not modify your `go.mod, go.sum` anymore
  - those two commands report if an module or checksum needs to be added or updated
  - use `go get` or `go mod tidy` to do tha

#### Other

- :tada: [`runtime/metrics`](https://tip.golang.org/pkg/runtime/metrics/) is a new package that provides a stable interface for accessing runtime metrics
  - generalizes [`runtime.ReadMemStats`](https://tip.golang.org/pkg/runtime/#ReadMemStats) and others
- [`GODEBUG=inittrace=1`](https://tip.golang.org/pkg/runtime/#hdr-Environment_Variables) prints execution time and memory allocation for all `init` calls, e.g. to analyze startup time regressions
- See [State of Go @ FOSDEM 2021](https://www.youtube.com/watch?v=pNd_BM0Tg4E) for another overview of the additions and changes coming with this release
- :racing_car: the [linker is faster](https://tip.golang.org/doc/go1.16#linker), requires less memory and binaries are getting smaller as a result of more aggressive symbol pruning¹
- ¹ my naïve explanation of symbol pruning:
  - a symbol contains metadata about the addresses of its variables and functions ([source](http://nickdesaulniers.github.io/blog/2016/08/13/object-files-and-symbols/))
  - :bulb: you can show the symbols (symbol-table) of a binary using `nm <binary-file>` (running `strip <binary-file>` will remove them all)
  - _symbol pruning_ refers to the removal of symbols that are unused, e.g. uncalled functions from imported libraries

---

## A deeper look on the new `embed` package

- [initial pull request](https://go-review.googlesource.com/c/go/+/243942) and [go command support](https://go-review.googlesource.com/c/go/+/243945)
- [embed.go](https://tip.golang.org/src/embed/embed.go) is only 423 lines including comments and blank lines
- strings are stored in plain-text in the binary
- embedded files must be located inside the package (that imports `embed`) directory or its subdirectories
- you can embed files as `string`, `[]byte` or `embed.FS`
  - use `embed.FS` when importing multiple (a tree of) files into a variable
- using the embed directive requires to import `embed` even when using `string` or `[]byte` variables (use blank import `import _ "embed")
- `// go:embed pattern [pattern...]` where paths use `path.Match` patterns
  - `string` and `[]byte` can only have a single pattern
- file patterns can be quoted like strings to embed paths containing spaces
- directories are embedded recursively while paths beginning with _ or . are excluded
  - note that `// go:embed dir/*` will embed `dir/.foo` but `// go:embed dir` will not!
- patterns must not begin with ., .. or /
- invalid patters and matches (e.g. a symbolic link) will cause the compilation to fail
- embed directive can only be used with exported or unexported global variables at package scope
- run example: `source .env && go run ./cmd/embed :12345`
- calling `strings` on the binary shows that content is embedded as plain-text

## `os` and `io/fs`

- with `os.ReadDir`/`fs.ReadDir` we get to new functions for listing directory contents
  - those functions are a lot more efficient than `ioutil.ReadDir` because they do not require to call `stat` for each directory entry
  - for details see [Ben Hoyts' article](https://benhoyt.com/writings/go-readdir/)
  - fun fact, those functions are inspired by [Python's `os.scandir`](https://docs.python.org/3/library/os.html#os.scandir)
  - to get the performance benefit of `ReadDir` when walking directory trees you need to use [`filepath.WalkDir`](https://tip.golang.org/pkg/path/filepath/#WalkDir) instead of `filepath.Walk`, `filepath.WalkDir` passes `fs.DirEntry` instead of `os.FileInfo` into the `WalkFn`

```go
type FS interface {
    Open(name string) (File, error)
}
```

- [`io.FS`](https://tip.golang.org/pkg/io/fs/#FS) filesystem abstraction interface
  - prior art is [`afero`](https://github.com/spf13/afero) but it has a [much larger interface](https://pkg.go.dev/github.com/spf13/afero#Fs) (and functionality, e.g. supports writes)
  - `embed.FS` implements this so you can easily work with embedded directory trees
  - pretty useful for tests, e.g. to prevent side effects → [`testing/fstest.MapFS`](https://tip.golang.org/pkg/testing/fstest/#MapFS) is an in-memory file system
  - accept `fs.FS` and make your API filesystem agnostic, i.e. work transparently with a local directory, s3 block storage, remote file system etc.
  - :bulb: you can easily test your custom file systems by passing them into [`fstest.TestFS`](https://tip.golang.org/pkg/testing/fstest/#TestFS)

```go
type DeleteFS interface {
  FS

  Delete(name string) error
}
```

- if you need more functionality use interface composition to build a larger interface, take [ReadDirFS](https://tip.golang.org/pkg/io/fs/#ReadDirFS) or [`StatFS`](https://tip.golang.org/pkg/io/fs/#ReadDirFS) as an example

## Generics

[The proposal got accepted](https://github.com/golang/go/issues/43651#issuecomment-776944155) :tada: 

- depending on your personal opinion this might be good or bad news
- we will not see generics in a 1.x release
- [Ardan Labs three part blog article series about generics](https://www.ardanlabs.com/blog/2020/07/generics-01-basic-syntax.html)
- (since some time already) there are two ways to try out generics:
  1. compile the `go2go` translation tool from the `dev.go2go` branch of the [Go repo](https://go.googlesource.com/go) ([instructions](https://go.googlesource.com/go/+/refs/heads/dev.go2go/README.go2go.md))
  2. use the [go2go playgound](https://go2goplay.golang.org/)

```go
func print[T any](value T) {
	fmt.Printf("type=%T\tvalue=%#v\n", value, value)
}

print([]int{1, 2, 3, 4})
print("Hello, Golang Leipzig!")
```

[playground](https://go2goplay.golang.org/p/cGQWVpzboRI)

- `[T any]` is the type parameter list
  - `any` is a type [_constraint_](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md#constraints) (limiting the set of allowed types that can be used for `T`)
  - `any` means "unconstrained" or allow any type for `T`

```go
type Number interface {
	type int, float64
}

func sum[T, S Number](a T, b S) S {
	return S(a) + b
}

sum(1, 3.1415)
```

[playground](https://go2goplay.golang.org/p/JWm9mDpsvVf)

- a type parameter list can contain multiple type identifiers
- ... we will talk likely have a separate talk that covers Generics in depth

## Questions :question: