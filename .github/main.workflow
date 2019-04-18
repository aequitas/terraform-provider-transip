workflow "Release" {
  on = "push"

  resolves = [
    "Build",
    "Test",
  ]
}

action "Build" {
  uses = "docker://golang:1.11"
  runs = ["sh", "-c", "go get -d && go build"]
}

action "Test" {
  uses = "docker://golang:1.11"
  runs = ["sh", "-c", "go get -d && go test"]
}

# action "Upload release Darwin" {
#   uses = "JasonEtco/upload-to-release@master"
#   args = "terraform-provider-transip_*_darwin_amd64.tgz"
#   secrets = ["GITHUB_TOKEN"]
#   needs = ["Build releases"]
# }


# action "Upload release Linux" {
#   uses = "JasonEtco/upload-to-release@master"
#   args = "terraform-provider-transip_*_linux_amd64.tgz"
#   secrets = ["GITHUB_TOKEN"]
#   needs = ["Build releases"]
# }
