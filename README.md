# Terraform Registry implementation (name to be determined)

## Description

This is a very simple, minimal implementation of the terraform registry API. It doesn't implement all methods, nor is it tested in anyway under stress. Use at your own risk!


## How to build

Just run `go build` in the root of the project, provided your `GOPATH` and `GOROOT` are setup properly.


## How to run

The executable will run on port 1234, http (no TLS). You will have to provide some commandline flags to make it work for your situation.

| Parameter    | Description                    | Default |
|  ----------- | ------------------------------ | ------- |
| -module_path | Base path for module storage   |         |
| -port        | Port to listen on              | 1234    |
| -tls_cert    | path to the server certificate |         |
| -tls_key     | path to the server key         |         |

When `tls_cert` AND `tls_key` are set the server will use https. Otherwise it defaults to http.
