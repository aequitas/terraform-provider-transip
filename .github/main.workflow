workflow "New workflow" {
  on = "push"
  resolves = ["cedrickring/golang-action@1.1.0"]
}

action "cedrickring/golang-action@1.1.0" {
  uses = "cedrickring/golang-action@1.1.0"
}
