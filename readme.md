# Automation Deployment

A Go REST API automatically built, tested, and deployed to AWS EC2 via GitHub Actions. Docker images are stored in Amazon ECR and pulled by the EC2 instance at deploy time.

## Architecture

```
GitHub Actions
    │
    ├── CI: build + test Go app
    │
    └── CD:
        ├── Build Docker image
        ├── Push to Amazon ECR
        └── SSH into EC2 → pull image → run container
```

**Infrastructure (Terraform):**
- VPC with public subnet
- EC2 instance (Amazon Linux 2023) with IAM role to pull from ECR
- Amazon ECR repository

**Application:**
- Go 1.25 + Gin framework
- Endpoints: `GET /` and `GET /users`

---

## Prerequisites

Make sure you have these installed:

- [Go 1.25+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Terraform](https://developer.hashicorp.com/terraform/install)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)

---

## Run Locally

### 1. Clone the repo

```bash
git clone https://github.com/<your-username>/automation-deployment.git
cd automation-deployment
```

### 2. Start the app

```bash
cd src
go mod download
go run ./cmd
```

App will be available at `http://localhost:8080`.

### 3. Run tests

```bash
cd src
go test ./test/...
```

---

## Deploy to AWS

### 1. Set up AWS credentials

Create an IAM user with the following permissions and generate an access key:

```json
{
  "Effect": "Allow",
  "Action": [
    "ec2:*", "ecr:*",
    "iam:CreateRole", "iam:DeleteRole", "iam:GetRole", "iam:PassRole",
    "iam:AttachRolePolicy", "iam:DetachRolePolicy",
    "iam:ListAttachedRolePolicies", "iam:ListRolePolicies",
    "iam:CreateInstanceProfile", "iam:DeleteInstanceProfile",
    "iam:GetInstanceProfile", "iam:AddRoleToInstanceProfile",
    "iam:RemoveRoleFromInstanceProfile", "iam:ListInstanceProfilesForRole"
  ],
  "Resource": "*"
}
```

Configure the AWS CLI:

```bash
aws configure --profile devops-automation
```

### 2. Create an EC2 key pair

```bash
aws ec2 create-key-pair \
  --key-name automation-deployment-key \
  --query 'KeyMaterial' \
  --output text > automation-deployment-key.pem

chmod 400 automation-deployment-key.pem
```

### 3. Provision infrastructure with Terraform

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:

```hcl
aws_region          = "ap-southeast-1"
project_name        = "automation-deployment"
ecr_repository_name = "automation-deployment"
key_name            = "automation-deployment-key"
instance_type       = "t3.micro"
```

Apply:

```bash
export AWS_PROFILE=devops-automation
terraform init
terraform apply
```

Note the outputs — you will need them for the next step:

```bash
terraform output ecr_repository_url
terraform output ec2_public_ip
terraform output app_url
```

### 4. Set up GitHub Actions secrets

Go to your GitHub repo → **Settings → Secrets and variables → Actions** and add:

| Secret | Value |
|--------|-------|
| `AWS_ACCESS_KEY_ID` | Your IAM user access key ID |
| `AWS_SECRET_ACCESS_KEY` | Your IAM user secret access key |
| `AWS_REGION` | e.g. `ap-southeast-1` |
| `EC2_HOST` | EC2 public IP from `terraform output ec2_public_ip` |
| `EC2_SSH_KEY` | Full contents of `automation-deployment-key.pem` |

To copy the key on macOS:

```bash
cat automation-deployment-key.pem | pbcopy
```

### 5. Trigger the pipeline

Push to `master` or `main` to trigger the CI/CD pipeline:

```bash
git push origin master
```

GitHub Actions will:
1. Build and test the Go app
2. Build the Docker image and push to ECR
3. SSH into EC2, pull the image, and run the container

### 6. Access the app

```bash
terraform output app_url
# e.g. http://54.151.xxx.xxx:8080
```

Open the URL in your browser or test with curl:

```bash
curl http://<ec2_public_ip>:8080
curl http://<ec2_public_ip>:8080/users
```

---

## Project Structure

```
.
├── .github/
│   └── workflows/
│       └── ci-cd.yml        # GitHub Actions CI/CD pipeline
├── src/
│   ├── cmd/
│   │   └── main.go          # App entrypoint
│   ├── internal/
│   │   ├── handler/         # HTTP handlers
│   │   └── model/           # Data models
│   ├── test/                # Tests
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
└── terraform/
    ├── main.tf              # Root module
    ├── variables.tf
    ├── outputs.tf
    ├── terraform.tfvars.example
    └── modules/
        ├── ec2/             # EC2 instance + IAM role
        ├── ecr/             # ECR repository
        └── networking/      # VPC, subnet, security group
```

---

## Cleaning Up

To destroy all AWS resources and avoid charges:

```bash
cd terraform
terraform destroy
```
