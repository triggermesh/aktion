workflow "knative test" {
  on = "push"
  resolves = [
    "First Action",
  ] 
}

action "First Action" {
  uses = "cab105/test-repo@master"
  runs = "echo"
  env = {
    FOO = "BAR"
  }
  args = "Hello world $FOO"
}
