workflow "hello multi-action test" {
  on = "push"
  resolves = [
    "First Action",
  ] 
}

action "First Action" {
  uses = "docker://centos"
  needs = "Second Action"
  runs = "echo"
  env = {
    FOO = "BAR"
  }
  args = "Hello world"
}

action "Second Action" {
  uses = "docker://centos"
  runs = "echo"
  args = [
    "tekton",
    "pipeline"
  ]
}
