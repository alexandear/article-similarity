# Article Similarity

HTTP server to store and search similar articles.

## Getting started

Run server and storage containers with Compose:

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

## Scalability

See [SCALEME](SCALEME.md) file.

## Technologies

There are HTTP server written on Golang and `mongodb` storage.

## Development

> Prerequisites: `docker`, `docker-compose`, `go@1.15`, `make` must be installed.

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

There are unit and integration tests. Unit tests placed in `_test.go` files,
end-to-end in `test` directory.

Run unit tests:

```shell
make test
```

End-to-end test suite builds server from sources, runs `docker-compose up` and perform requests to server container.
It can be executed:

```shell
make test-it
```

### Docker

Build docker image `article-similarity:latest`:

```shell
make docker
```

Build, run linter and tests in dev docker image `article-similarity-dev:latest`:

```shell
make docker-dev
```

### CI

There are configured GitHub actions for build, lint, run unit and integration tests. 
See [.github/workflows](.github/workflows) directory.
