FROM golang:alpine as build

RUN mkdir /registry
ADD . /src/github.com/erikvanbrakel/anthology

WORKDIR /src/github.com/erikvanbrakel/anthology

ENV GOPATH /

RUN go build && cp ./anthology /registry/anthology

FROM alpine:latest

COPY --from=build /src/github.com/erikvanbrakel/anthology/anthology /registry/anthology
COPY --from=build /src/github.com/erikvanbrakel/anthology/test/modules /modules

WORKDIR /registry

EXPOSE 8082

CMD ["-port=8082","-module_path=/modules"]
ENTRYPOINT ["./anthology"]
