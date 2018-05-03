FROM golang:alpine as build

RUN mkdir /registry
ADD . /src/github.com/erikvanbrakel/anthology

WORKDIR /src/github.com/erikvanbrakel/anthology

ENV GOPATH /

RUN go build && cp ./anthology /registry/anthology

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=build /src/github.com/erikvanbrakel/anthology/anthology /registry/anthology

WORKDIR /registry

EXPOSE 8082

CMD ["--port=8082","--filesystem.basepath=/modules","--backend=filesystem"]
ENTRYPOINT ["./anthology"]
