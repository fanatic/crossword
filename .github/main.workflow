workflow "Build/Deploy" {
  on = "push"
  resolves = [
    "Frontend (S3)",
    "Backend (Lambda)",
  ]
}

action "Frontend (S3)" {
  uses = "./.github/frontend"
  args = ["cd client ; yarn install && yarn deploy"]
  secrets = ["AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"]
}

action "Backend (Lambda)" {
  uses = "./.github/backend"
  args = ["cd server; ./deploy.sh"]
  secrets = ["AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"]
}
