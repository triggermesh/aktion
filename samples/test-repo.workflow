workflow "github repo test" {
  on = "push"
  resolves = [
    "First Action",
  ] 
}

action "First Action" {
  uses = "triggermesh/aktion/samples/images@master"
  runs = "echo"
  env = {
    FOO = "BAR"
  }
  args = "Hello world $FOO"
}
