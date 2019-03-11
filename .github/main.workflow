workflow "Release" {
  on = "release"
  resolves = [
    "Upload release Darwin",
    "Upload release Linux",
  ]
}

action "Build releases" {
  uses = "cedrickring/golang-action@1.1.0"
}

action "Upload release Darwin" {
  uses = "JasonEtco/upload-to-release@master"
  args = "terraform-provider-transip_*_darwin_amd64.tgz"
  secrets = ["GITHUB_TOKEN"]
  needs = ["Build releases"]
}

action "Upload release Linux" {
  uses = "JasonEtco/upload-to-release@master"
  args = "terraform-provider-transip_*_linux_amd64.tgz"
  secrets = ["GITHUB_TOKEN"]
  needs = ["Build releases"]
}
