# A quick tour of Go 1.16 @ 16th Golang Leipzig Meetup

Our 16th Meetup at the 16th of February with a quick tour on Go 1.16 :blush: .
We will also have an outlook on the Go generics proposal.

## `embed`

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

### Open Questions

- io/fs while being read-only, how to easily add write interfaces?
- generics
- quick note about the "deprectation" of io/ioutil
