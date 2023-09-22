# `todo-list-server` provider for Terraform

## `terraform-provider-todo`

This Terraform provider code for the `todo-list-server` was migrated here from [spkane/todo-for-terraform](https://github.com/spkane/todo-for-terraform/) to make it easier to deploy this to Hashicorp's Terraform Provider Registry.

The most recent release of the `todo-list-server` can be found [here](https://github.com/spkane/todo-for-terraform/releases).

### Build provider

Run the following command to build the provider

```shell
go build -o terraform-provider-todo
```

### Test provider

```shell
make test
make testacc
```

- If you want to test out the provider with the `terraform` CLI.

Edit `$HOME/.terraformrc`

```hcl
provider_installation {
  dev_overrides {
    "spkane/todo" = "/home/me/go/path/bin/"
  }
  direct {}
}
```

And then create and test a few runs based on the files under examples.

### Release

- In a PR, consider updating the Go version and modules.

```sh
go get -u ./...
go mod tidy
```

- In a PR, consider updating the GitHub Actions.

```sh
vi .github/workflows/release.yaml
```

- Rebuild, test, and merge.

- Add a new version tag to the `main` branch.
  - This will trigger the GitHub Actions that build and push the new release.

### Documentation

Documentation is generated with [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) and exists in the [docs](./docs/) directory.

## Pre-Commit Hooks

- See: [pre-commit](https://pre-commit.com/)
  - [pre-commit/pre-commit-hooks](https://github.com/pre-commit/pre-commit-hooks)
  - [antonbabenko/pre-commit-terraform](https://github.com/antonbabenko/pre-commit-terraform)

### Install

#### Local Install (macOS)

- **IMPORTANT**: All developers committing any code to this repo, should have these pre-commit hooks installed locally. Github actions may also run these at some point, but it is generally faster and easier to run them locally, in most cases.

```sh
brew install pre-commit jq shellcheck shfmt git-secrets go-critic golangci-lint
go install github.com/BurntSushi/toml/cmd/tomlv@master
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install golang.org/x/tools/cmd/goimports@latest

mkdir -p ${HOME}/.git-template/hooks
git config --global init.templateDir ${HOME}/.git-template
```

- Close and reopen your terminal
- Make sure that you run these commands from the root of this git repo!

```sh
cd terraform-provider-todo
pre-commit init-templatedir -t pre-commit ${HOME}/.git-template
pre-commit install
```

- Test it

```sh
pre-commit run -a
git diff
```

### Checks

See:

- [.pre-commit-config.yaml](./.pre-commit-config.yaml)

#### Configuring Hooks

- [.pre-commit-config.yaml](./.pre-commit-config.yaml)

## TODO (FIXME)

- Resource READ should handle 404s that might indicate that a resource has been deleted out of band.
