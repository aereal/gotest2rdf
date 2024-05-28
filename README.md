[![status][ci-status-badge]][ci-status]
[![PkgGoDev][pkg-go-dev-badge]][pkg-go-dev]

# gotest2rdf

**gotest2rdf** transforms `go test -json` outputs into [Reviewdog Diagnostic Format][rdf].

It is useful with using [reviewdog][].

![the image with reviewdog -reporter github-pr-review][image]

## Synopsis

```sh
go test -json ./... | gotest2rdf | reviewdog -f rdjsonl
```

## Installation

See releases page or:

```sh
go install github.com/aereal/gotest2rdf@latest
```

## See also

- 
- [cmd/test2json]: the details of `go test -json`.

## License

See LICENSE file.

[pkg-go-dev]: https://pkg.go.dev/github.com/aereal/gotest2rdf
[pkg-go-dev-badge]: https://pkg.go.dev/badge/aereal/gotest2rdf
[ci-status-badge]: https://github.com/aereal/gotest2rdf/workflows/CI/badge.svg?branch=main
[ci-status]: https://github.com/aereal/gotest2rdf/actions/workflows/CI
[rdf]: https://github.com/reviewdog/reviewdog/tree/master/proto/rdf
[test2json]: https://pkg.go.dev/cmd/test2json
[reviewdog]: https://github.com/reviewdog/reviewdog
[image]: /docs/image.png
