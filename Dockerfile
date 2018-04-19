FROM golang:alpine as build

RUN mkdir /registry
ADD . /src/github.com/erikvanbrakel/terraform-registry

WORKDIR /src/github.com/erikvanbrakel/terraform-registry

ENV GOPATH /

RUN go build && cp ./terraform-registry /registry/terraform-registry

FROM alpine:latest

COPY --from=build /src/github.com/erikvanbrakel/terraform-registry/terraform-registry /registry/terraform-registry
COPY --from=build /src/github.com/erikvanbrakel/terraform-registry/test/modules /modules

WORKDIR /registry

EXPOSE 8082

CMD ["-port=8082","-module_path=/modules"]
ENTRYPOINT ["./terraform-registry"]
