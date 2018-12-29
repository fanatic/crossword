workflow "New workflow" {
  on = "push"
  resolves = ["GitHub Action for AWS", "GitHub Action for AWS-1"]
}

action "GitHub Action for AWS" {
  uses = "actions/aws/cli@8d31870"
  args = "sh -c \"cd client ; yarn deploy\""
  secrets = ["AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"]
}

action "GitHub Action for AWS-1" {
  uses = "actions/aws/cli@8d31870"
  args = "sh -c \"cd server; ./deploy.sh\""
}
