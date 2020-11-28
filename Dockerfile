FROM golang:1.15-alpine as build

RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache git make

WORKDIR ./src/github.com/devchallenge/article-similarity

COPY ./assets /assets
COPY . ./

RUN make build && cp ./bin/article-similarity /usr/local/bin/

FROM alpine

COPY --from=build /usr/local/bin/ /usr/local/bin/
COPY --from=build /assets ./assets

ENTRYPOINT ["article-similarity", "--host=0.0.0.0", "--port=80"]
