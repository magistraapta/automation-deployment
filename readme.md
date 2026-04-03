Step 3 — Deploy: Do This In Order

  A. AWS Prerequisites

  1. Create an EC2 key pair in the AWS console → save the .pem file locally
  2. Create an IAM user with AmazonEC2FullAccess, AmazonECRFullAccess, and AmazonVPCFullAccess →
  save its access key + secret

  B. Provision AWS Infrastructure with Terraform

  cd terraform
  cp terraform.tfvars.example terraform.tfvars
  # Edit terraform.tfvars — set key_name to your EC2 key pair name

  terraform init
  terraform plan
  terraform apply

  After apply, note the two outputs:
  - ecr_repository_url → e.g.
  123456789.dkr.ecr.ap-southeast-1.amazonaws.com/automation-deployment
  - ec2_public_ip → e.g. 54.123.45.67

  C. Add GitHub Secrets

  Go to Settings → Secrets and variables → Actions in your repo and add:

  ┌───────────────────────┬──────────────────────────────────────────┐
  │        Secret         │                  Value                   │
  ├───────────────────────┼──────────────────────────────────────────┤
  │ AWS_ACCESS_KEY_ID     │ IAM user access key                      │
  ├───────────────────────┼──────────────────────────────────────────┤
  │ AWS_SECRET_ACCESS_KEY │ IAM user secret key                      │
  ├───────────────────────┼──────────────────────────────────────────┤
  │ AWS_REGION            │ e.g. ap-southeast-1                      │
  ├───────────────────────┼──────────────────────────────────────────┤
  │ EC2_HOST              │ the ec2_public_ip from Terraform output  │
  ├───────────────────────┼──────────────────────────────────────────┤
  │ EC2_SSH_KEY           │ contents of your .pem file (entire text) │
  └───────────────────────┴──────────────────────────────────────────┘

  D. Push to master → Pipeline Runs

  git add .
  git commit -m "add terraform and fix ci-cd pipeline"
  git push origin dev
  # then merge dev → master to trigger the workflow

  ---
  How the Full Pipeline Works

  push to master
       │
       ▼
    [CI job]
    ├── setup Go 1.23
    ├── go build ./cmd/main.go
    └── go test ./test/...
       │ (only if CI passes)
       ▼
    [CD job]
    ├── configure AWS credentials
    ├── login to ECR
    ├── docker build + push to ECR (tagged with git SHA)
    └── SSH into EC2 → pull new image → restart container

  One thing to note: The go-version: 1.25.7 in the original workflow would have caused the
  pipeline to fail immediately — Go 1.23 is the latest stable release.