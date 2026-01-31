# quiz ECS デプロイ用 Makefile

# 変数
SERVICE_NAME = quiz
AWS_REGION ?= ap-northeast-1
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)
ECR_REGISTRY = $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_REPOSITORY = $(ECR_REGISTRY)/$(SERVICE_NAME)-$(ENV)
GIT_COMMIT_HASH ?= $(shell git rev-parse --short HEAD)
ENV ?= stg

# ローカルでビルドする
# e.g. make build
build:
	docker build \
		-t $(SERVICE_NAME):$(GIT_COMMIT_HASH) \
		-t $(SERVICE_NAME):latest \
		.

# ECRにログインする
# e.g. make ecr-login
ecr-login:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_REGISTRY)

# Dockerイメージをタグ付けしてECRにプッシュする
# e.g. make push ENV=stg
push: build ecr-login
	docker tag $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH)
	docker push $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH)
	docker tag $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(IMAGE_REPOSITORY):latest
	docker push $(IMAGE_REPOSITORY):latest

# ビルド→ECRへpushをまとめて行う
# e.g. make deploy ENV=stg
deploy: build push
	@echo "Deployed $(SERVICE_NAME) to ECR: $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH)"

# ローカルで実行する（Docker経由）
# e.g. make run-local
run-local:
	docker run --rm -p 8080:8080 \
		-e DB_HOST=$(DB_HOST) \
		-e DB_PORT=$(DB_PORT) \
		-e DB_USER=$(DB_USER) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e DB_NAME=$(DB_NAME) \
		-e LOG_LEVEL=$(LOG_LEVEL) \
		-e CONFIG_PATH=$(CONFIG_PATH) \
		$(SERVICE_NAME):latest

# ローカルで直接実行する（Go経由）
# e.g. make run
run:
	go run ./cmd/main.go

# Swaggerドキュメントを生成する
# e.g. make swagger
swagger:
	swag init -g cmd/main.go -o docs

# テストを実行する
# e.g. make test
test:
	go test -v ./...

# go mod tidyを実行する
# e.g. make tidy
tidy:
	go mod tidy

# クリーンアップ
# e.g. make clean
clean:
	docker rmi $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(SERVICE_NAME):latest || true
	docker rmi $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH) $(IMAGE_REPOSITORY):latest || true

.PHONY: build ecr-login push deploy run-local run swagger test tidy clean
