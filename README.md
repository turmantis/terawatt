# Terawatt
A simple, configuration free, version manager for terraform. 

## Install
```shell
go install github.com/turmantis/terawatt@latest
```

## Example
Terawatt is a simple wrapper around terraform. All args are passed through to terraform.
```shell
$ echo 'terraform {required_version = "~> 1.2.3"}' > foo.tf
$ terawatt plan
Downloading: https://releases.hashicorp.com/terraform/1.2.3/terraform_1.2.3_linux_amd64.zip

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## Version Resolution
1. Check if `.terraform-version` exists. If it does, stop and use that version.
2. Scan all files ending in `.tf` and looking `required_version` statement. If it exists, use that version.
3. Fetch the available versions of terraform and use the latest stable version.

Currently only exact version expression can be parsed from hcl files (eg `~> 1.2.3-foo`).

## Compare With
|                                                 | terawatt | tfenv |
|-------------------------------------------------|----------|-------|
| Automatically download terraform                | 🟢       | 🔴    |   
| Parses version from `required_version` from hcl | 🟢       | 🔴    |
| Parses version from `.terraform-version`        | 🟢       | 🟢    |
| Parses version expressions (eg `>=1.0<2`)       | 🔴       | 🟢    |