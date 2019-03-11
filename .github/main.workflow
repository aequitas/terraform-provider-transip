workflow "New workflow" {
  on = "push"
  resolves = ["Make"]
}

action "Make" {
  uses = "Make"
  runs = "make"
}
