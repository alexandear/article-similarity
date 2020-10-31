FROM golang:1.15-alpine as build

RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache git make

WORKDIR ./src/github.com/devchallenge/article-similarity

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build && cp ./bin/article-similarity /usr/local/bin/

FROM alpine

COPY --from=build /usr/local/bin/article-similarity /usr/local/bin/article-similarity

ENV HOST 0.0.0.0
ENV PORT 80

EXPOSE 80

ENTRYPOINT ["article-similarity"]
