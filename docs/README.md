# Shell Provider

`shell` is a special provider that exists to provide an interface between Terraform and external scripts.

Using this provider it is possible to write scripts that can participate in the Terraform workflow by implementing a specific protocol.

This provider is intended to be used for situations where implementing a terraform provider for an interface would be prohibitive. As we are working with scripts, trade-offs had to be made to ensure consistency when running each of these examples. 

## Example Usage

```
# Configure the Shell provider
provider "shell" {
  working_directory = "/tmp"

  variables = {
    USERNAME = "MyUser"
    PASSWORD = "Password1"
  }

  prune = [
    "Remove this prefix from outputs"
  ]
}

# Map to a script
resource "shell" "run_script_x" {
  # ...
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

- `working_directory` - (Optional) The direction where the scripts will run from. It will be defaulted to a temporary location.
- `variables` - (Optional) Environment variables that will be passed down to all scripts using this provider.
- `prune` - (Optional) A workaround to prune text that cannot be suppressed from outputs.
