workflow "secrets test" {
  on = "push"
  resolves = [
    "First Action",
  ] 
}

action "First Action" {
  uses = "docker://centos"
  runs = "echo"
  args = "Hello world"
  secrets = ["BAR", "BAZ"]
}
