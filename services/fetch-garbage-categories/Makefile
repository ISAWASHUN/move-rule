# fetch-garbage-categories Lambda デプロイ用 Makefile

# 変数
SERVICE_NAME = fetch-garbage-categories
AWS_REGION ?= ap-northeast-1
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)
ECR_REGISTRY = $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_REPOSITORY = $(ECR_REGISTRY)/$(SERVICE_NAME)
GIT_COMMIT_HASH ?= $(shell git rev-parse --short HEAD)
ENV ?= stg

# Lambda関数名
LAMBDA_FUNCTION_NAME = $(SERVICE_NAME)-$(ENV)

# ローカルでビルドする
# e.g. make build
build:
	docker build -t $(SERVICE_NAME):$(GIT_COMMIT_HASH) .
	docker tag $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(SERVICE_NAME):latest

# ECRにログインする
# e.g. make ecr-login
ecr-login:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_REGISTRY)
# Dockerイメージをタグ付けしてECRにプッシュする
# e.g. make push
push: ecr-login
	docker tag $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH)
	docker tag $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(IMAGE_REPOSITORY):latest
	docker push $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH)
	docker push $(IMAGE_REPOSITORY):latest

# Lambda関数のコードを更新する
# e.g. make update-function ENV=stg
update-function:
	aws lambda update-function-code \
		--function-name $(LAMBDA_FUNCTION_NAME) \
		--image-uri $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH) \
		--region $(AWS_REGION)

# ビルド→ECRへpush→Lambda関数の更新をまとめて行う
# e.g. make deploy ENV=stg
deploy: build push update-function
	@echo "Deployed $(SERVICE_NAME) to $(LAMBDA_FUNCTION_NAME) with image $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH)"

# ローカルで実行する（Docker経由）
# e.g. make run-local S3_BUCKET=my-bucket S3_PREFIX=data
run-local:
	docker run --rm \
		-e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
		-e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
		-e AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN) \
		-e S3_BUCKET=$(S3_BUCKET) \
		-e S3_PREFIX=$(S3_PREFIX) \
		-e AWS_REGION=$(AWS_REGION) \
		$(SERVICE_NAME):latest

# ローカルで直接実行する（Go経由）
# e.g. make run
run:
	go run ./cmd/main.go

# テストを実行する
# e.g. make test
test:
	go test -v ./...

# クリーンアップ
# e.g. make clean
clean:
	docker rmi $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(SERVICE_NAME):latest || true
	docker rmi $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH) $(IMAGE_REPOSITORY):latest || true

.PHONY: build ecr-login create-ecr-repo push update-function deploy run-local run test clean
