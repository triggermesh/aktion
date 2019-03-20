workflow "knative test" {
  on = "push"
  resolves = [
    "First Action",
  ] 
}

action "First Action" {
  uses = "docker://centos"
  needs = "With Secrets"
  runs = "echo"
  env = {
    FOO = "BAR"
  }
  args = "Hello world"
  secrets = ["BAR", "BAZ"]
}

action "With Secrets" {
  uses = "docker://centos"
  runs = "echo"
  secrets = ["FOO", "BAR"]
}
