workflow "knative test" {
  on = "push"
  resolves = [
    "First Action",
  ] 
}

action "First Action" {
  uses = "docker://centos"
  runs = "echo"
  env = {
    FOO = "BAR"
  }
  args = "Hello world $FOO"
}
