provider "shell" {
  variables = {
    Message = "Hello, World!"
  }
}

resource "shell" "default" {
  create = ["pwsh", "-File", "${path.module}/Command.ps1", "-Action", "create"]
  read   = ["pwsh", "-File", "${path.module}/Command.ps1", "-Action", "read"]
  update = ["pwsh", "-File", "${path.module}/Command.ps1", "-Action", "update"]
  delete = ["pwsh", "-File", "${path.module}/Command.ps1", "-Action", "delete"]
  query = {
    description = "This should fail to run"
  }
}

output "id" { value = shell.default.id }
