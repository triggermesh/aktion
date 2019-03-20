workflow "foo" {
	resolves = "aaa"
}

action "aaa" {
	needs = "bbb"
}

action "bbb" {
	needs = ["aaa", "ccc"]
}
