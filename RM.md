# SuperOrbital Infrastructure

This repo contains our infrastructure configuration.  Each directory is a separate environment.

## First things first

Secrets have been migrated to [Doppler](https://dashboard.doppler.com/workplace/fa9c535d8ebb1789db28/projects/infrastructure).

- [Install doppler](https://docs.doppler.com/docs/enclave-installation) (most likely via `brew install dopplerhq/cli/doppler`)
- Authenticate locally with `doppler login` and your `@superorbital.io` email address. Tammer can give you expanded rights once you have logged on initially.
- Configure doppler for this project by running `doppler setup -p infrastructure -c prd`

The `bin/run` script will grab your Doppler token and provide it to the terraform commands.

We use git-lfs to vendor binaries.  We try to vendor based on platform (`bin/Linux/*`), but we haven't done a great job of remembering to vendor Darwin bins.  :shrug:

We use `aws-vault` for AWS authentication.  Any AWS terraform directories should have an `aws-vault-profile-name` file that tells `./bin/run` which profile to use.  This only works if we all use the same profile names, which is sub-optimal and why we started developing [Cludo](https://github.com/superorbital/cludo).  :shrug:

## Pre-Commit Hooks

- See: [pre-commit](https://pre-commit.com/)
  - [pre-commit/pre-commit-hooks](https://github.com/pre-commit/pre-commit-hooks)
  - [antonbabenko/pre-commit-terraform](https://github.com/antonbabenko/pre-commit-terraform)

### Install

#### Local Install (macOS)

- **IMPORTANT**: All developers committing any code to this repo, should have these pre-commit hooks installed locally. Github actions may also run these at some point, but it is generally faster and easier to run them locally, in most cases.

```sh
brew install pre-commit terraform-docs tfenv tflint tfsec checkov terrascan infracost tfupdate minamijoyo/hcledit/hcledit jq shellcheck shfmt git-secrets

mkdir -p ${HOME}/.git-template/hooks
git config --global init.templateDir ${HOME}/.git-template
```

- Close and reopen your terminal
- Make sure that you run these commands from the root of this git repo!

```sh
cd infrastructure
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
- [.tflint.hcl](./.tflint.hcl)

#### Currently Disabled or Limited Checks

- [tfupdate](https://github.com/minamijoyo/tfupdate) - Update required versions to stay current.
- [terraform_validate](https://www.terraform.io/cli/commands/validate) - This validates the configuration files in a directory, referring only to the configuration and not accessing any remote services such as remote state, provider APIs, etc.
- [infracost](https://www.infracost.io/) - Predicted infrastructure costs
  - Requires [a free API key](https://www.infracost.io/docs/#2-get-api-key)
- [checkov](https://www.checkov.io/) - Configuration security scanning (Terraform, Kubernetes, etc.)
  - Noisy until some time is spent improving the files and configuring the tool for the things that we care about.
- [tfsec](https://github.com/aquasecurity/tfsec) - Static analysis to spot potential misconfigurations.
  - `tfsec` is currently disabled. Like checkov this will require some time fixing existing terraform and fine tuning the rules.
- [terrascan](https://runterrascan.io/) - Detect compliance and security violations across Infrastructure as Code (IaC)
  - Noisy until some time is spent improving the files and configuring the tool for the things that we care about.
