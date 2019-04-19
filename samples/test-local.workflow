workflow "local repo test" {
  on = "push"
  resolves = [
    "First Action",
  ]
}

action "First Action" {
  uses = "./samples/test-images"
  runs = "echo"
  env = {
    FOO = "BAR"
  }
  args = "Hello world $FOO"
}
