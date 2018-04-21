# Anthology, a private Terraform Registry

## Description

Anthology is a reimplementation of the Terraform Registry API, intended to be used when your modules can't, shouldn't or don't
need to be public. For all means and purposes it works in the same way as the [public registry][terraform-registry].


## How to use

### Using Docker

Every release is automatically published to the [Docker Hub][docker-hub]. You can set commandline parameters by overriding the
command.

__running on port `80`, using `my-module-bucket` for storage:__

`docker run -p 80:80 erikvanbrakel/anthology -port=80 -bucket=my-module-bucket`

__using docker-compose__
```yaml
version: '2.1'

services:

  registry:
    command: -port=80 -bucket=my-module-bucket
    build: erikvanbrakel/anthology:latest
    ports:
      - 80:80
```

### AWS + terraform

The easiest way to deploy is to use the [anthology module][anthology-module] in the [public registry][terraform-registry].

```hcl
module "anthology" {
  source  = "erikvanbrakel/anthology/aws"
  version = "0.0.2"

  storage_bucket = "this-bucket-stores-my-modules"
  tld            = "example.com"                   # the registry will be hosted at registry.example.com
}

```

__WARNING WARNING WARNING__

This module provisions several resources, among which compute and storage components. This is not free, so make sure you are
aware of the cost before provisioning!


## Command line parameters

| Parameter    | Description                       | Default |
|  ----------- | --------------------------------- | ------- |
| -module_path | Base path for module storage      |         |
| -port        | Port to listen on                 | 1234    |
| -tls_cert    | path to the server certificate    |         |
| -tls_key     | path to the server key            |         |
| -bucket      | name of the S3 bucket for storage |         |


When `tls_cert` AND `tls_key` are set the server will use https. Otherwise it defaults to http.

[terraform-registry]: https://registry.terraform.io/
[anthology-module]: https://registry.terraform.io/modules/erikvanbrakel/anthology/aws/
[docker-hub]: https://hub.docker.com/r/erikvanbrakel/anthology/
