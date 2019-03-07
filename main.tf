provider "shell" {
  working_directory = "/"
}

resource "shell" "experiment" {
  #   program = ["bash", "-c", "echo", "{\"hello\":\"world\"}"]
  create = ["bash", "script.sh"]
  update = ["bash", "script.sh"]
  read   = ["bash", "script.sh"]
  delete = ["bash", "script.sh"]

  query {
    hello = "world"
  }
}
