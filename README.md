[![Gitter chat](https://badges.gitter.im/anthology-registry/community.png)](https://gitter.im/anthology-registry/community)

# Anthology, a private Terraform Registry

## Description

Anthology is a reimplementation of the Terraform Registry API, intended to be used when your modules can't, shouldn't
or don't need to be public. For all means and purposes it works in the same way as the [public registry][terraform-registry].


## How to use

### Using Docker

Every release is automatically published to the [Docker Hub][docker-hub]. You can set commandline parameters by
overriding the command.

__running on port `80`, using `my-module-bucket` for storage:__

`docker run -p 80:80 erikvanbrakel/anthology --port=80 --backend=s3 --s3.bucket=my-module-bucket`

__using docker-compose__
```yaml
version: '2.1'

services:

  registry:
    command: --port=80 --backend=s3 --s3.bucket=my-module-bucket
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

This module provisions several resources, among which compute and storage components. This is not free, so make sure you
are aware of the cost before provisioning!


## Command line parameters

### Common parameters
| Parameter             | Description                       | Allowed                  | Default |
| --------------------- | --------------------------------- | ------------------------ | ------- |
| --port                | Port to listen on                 | 1-65535                  | 1234    |
| --backend             | Backend to use.                   | [memory, filesystem, s3] |         |
| --ssl.certificate     | Path to the server certificate    | Any valid path           |         |
| --ssl.key             | Path to the server certificate    | Any valid path           |         |

### Filesystem backend
| Parameter             | Description                       | Allowed                  | Default |
| --------------------- | --------------------------------- | ------------------------ | ------- |
| --filesystem.basepath | Base path for module storage      | Any valid path           |         |

### S3 backend
| Parameter             | Description                       | Allowed                    | Default |
| --------------------- | --------------------------------- | -------------------------- | ------- |
| --s3.bucket           | Name of the S3 bucket for storage | Any valid s3 bucket name   |         |
| --s3.endpoint         | Alternative S3 endpoint           | http[s]://[hostname]:[port]|         |

[terraform-registry]: https://registry.terraform.io/
[anthology-module]: https://registry.terraform.io/modules/erikvanbrakel/anthology/aws/
[docker-hub]: https://hub.docker.com/r/erikvanbrakel/anthology/
