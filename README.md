# Terawatt
A simple, configuration free, version manager for terraform. 

## Example
Terawatt is a simple wrapper around terraform. All args are passed through to terraform.
```shell
go install github.com/turmantis/terawatt@latest
echo 'terraform {required_version = "~> 1.2.3"}' > foo.tf
terawatt plan
```

## Version Resolution
If a `.terraform-version` file exists, that version will be used. Otherwise, all of the hcl files
in the current directory (ending in `.tf`) will be searched for the `required_version` statement. If
no required version is specified the latest stable version will be used.

Currently only exact version expression can be parsed from hcl `~> 1.2.3-foo`.
