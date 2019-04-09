workflow "local repo test" {
  on = "push"
  resolves = [
    "First Action",
  ]
}

action "First Action" {
  uses = "./docker-action"
  runs = "echo"
  env = {
    FOO = "BAR"
  }
  args = "Hello world $FOO"
}
