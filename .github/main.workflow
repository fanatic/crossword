workflow "Build/Deploy" {
  on = "push"
  resolves = [
    "Frontend (S3)",
    "Backend (Lambda)",
  ]
}

action "Frontend (S3)" {
  uses = "actions/aws/cli@8d31870"
  args = "sh -c \"cd client ; yarn deploy\""
  secrets = ["AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"]
}

action "Backend (Lambda)" {
  uses = "actions/aws/cli@8d31870"
  args = "sh -c \"cd server; ./deploy.sh\""
  secrets = ["AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"]
}

workflow "New workflow" {
  on = "push"
}
