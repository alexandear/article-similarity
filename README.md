# Article Similarity

HTTP server to store and search similar articles.

## Getting started

Build and run server with Compose:

```shell
docker-compose up
```

API is accessible via `http://localhost:80/`.

## API docs

API's description is in the [docs/api.md](./docs/api.md) file.

Additionally, server serves HTML documentation. Run `docker-compose up` and visit http://localhost:80/docs.

## Similarity algorithm

To find similarity between the content of articles used Levenshtein algorithm for words. Before Levenshtein algorithm is 
applied content preprocessing:
- remove articles `a, an, the` and punctuation `.,!?-`;
- content separated to word via whitespace characters ` \t\n\r`;
- text is lower-cased.

Algorithm works for English content only.

## Development

> Prerequisites: install `go@1.15` compiler and `make`.

### Code style

Consistent code style enforced by `gofmt`  and `golangci-lint` linters.

Format code:

```shell
make format
```

Run linter:

```shell
make lint
```

### Tests

There are unit and end-to-end tests. Unit and integration tests placed in `_test.go` files,
end-to-end in `test` directory.

Run unit tests:

```shell
make test
```

Run end-to-end tests:

```shell
make test-it
```
