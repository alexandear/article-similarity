# Article Similarity

HTTP server to store and search similar articles.

## API docs

API's description is in the [docs/api.md](./docs/api.md) file.

Additionally, server provides docs in HTML format via http://localhost:80/docs.

## Algorithm

To find similarity between the content of articles used Levenshtein algorithm for words. Before Levenshtein algorithm is 
applied content preprocessing:
- remove articles `a, an, the` and punctuation `.,!?-`;
- content separated to word via whitespace characters ` \t\n\r`;
- text is lower-cased.

Algorithm works for English content only.

## Tests

There are unit, integration and end-to-end tests. Unit and integration tests placed in `_test.go` files, 
end-to-end in `test` directory.
