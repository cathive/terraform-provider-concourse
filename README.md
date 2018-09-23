# Terraform Provider for Concourse CI

## Maintainers

This provider plugin is maintained by:

* Benjamin P. Jung <headcr4sh@gmail.com>
* Nick Larsen <nick@aptiv.co.nz>

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

## Building the Provider

Clone repository *outside* of your GOPATH:

```sh
$ mkdir -p ~/Projects; cd ~/Projects
$ git clone git@github.com:cathive/terraform-provider-concourse
```

Enter the provider directory and build the provider

```sh
$ cd ~/Projects/terraform-provider-concourse
$ make build
```

## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin).
After placing it into your plugins directory,  run `terraform init` to initialize it.

Documentation for all data providers and resources can be found in the subfolder [/docs/](./docs)
right here in this repository.
