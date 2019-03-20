workflow "on push" {
	on = "push"
	resolves = "go-ci"
}

action "go-ci" {
	uses = "docker://golang:latest"
	runs = "./script/cibuild"
}
