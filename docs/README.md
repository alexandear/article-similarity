# API docs

Contains converted from swagger spec API documentation.

# How to regenerate

Install [widdershins](https://mermade.github.io/widdershins/ConvertingFilesBasicCLI.html):

```shell
npm install -g widdershins
```

Run:

```shell
widdershins ./api/spec.yaml --code true --search false --language_tabs --outfile ./docs/api.md
```

For more available options [see](https://github.com/Mermade/widdershins#options).
