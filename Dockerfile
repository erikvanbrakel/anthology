FROM golang:alpine

RUN mkdir /registry
ADD . /src/github.com/erikvanbrakel/terraform-registry

WORKDIR /src/github.com/erikvanbrakel/terraform-registry

ENV GOPATH /

RUN go build && cp ./terraform-registry /registry/terraform-registry

WORKDIR /registry

EXPOSE 8082

ENTRYPOINT ./terraform-registry -port=8082 -module_path /src/github.com/erikvanbrakel/terraform-registry/test/modules
