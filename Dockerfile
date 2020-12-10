FROM golang:1.15-alpine as build

WORKDIR ./src

COPY ./assets /assets
COPY . ./

RUN go build -mod=vendor -o=./bin/article-similarity main.go && \
    cp ./bin/article-similarity /usr/local/bin/ && \
    rm -rf /go/src

FROM alpine

COPY --from=build /usr/local/bin/ /usr/local/bin/
COPY --from=build /assets ./assets

ENTRYPOINT ["article-similarity", "--host=0.0.0.0", "--port=80"]
