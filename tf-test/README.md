# Hands-on Testing

## Getting Started

* modify `~/.terraformrc`

```hcl
provider_installation {
  dev_overrides {
    "spkane/todo" = "/Users/me/dev/go/path/bin/"
  }
  direct {}
}
```

* Build the binary, by running `make` in the root of the git repo.
* Spin up a copy of the inventory service. `docker run --name todo-list --rm -p 8080:80 spkane/todo-list-server:latest`
* Run `terraform init`
* Run whatever other terraform commands you want.

## Cleaning up

* Stop the inventory service. `docker container stop todo-list`
* Comment out the `dev_overrides` in `~/.terraformrc`.
