ROOT_DIR = ../../

# ECR デプロイ用変数
SERVICE_NAME = db-migrator
AWS_REGION ?= ap-northeast-1
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)
ECR_REGISTRY = $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_REPOSITORY = $(ECR_REGISTRY)/$(SERVICE_NAME)-$(ENV)
GIT_COMMIT_HASH ?= $(shell git rev-parse --short HEAD)
ENV ?= stg
SUBNET_ID ?= subnet-0f6136e99974839ad
SECURITY_GROUP_ID ?= sg-0882ca3ecc9e2cd2f
CLUSTER ?= garbage-category-rule-quiz-$(ENV)

# ================================================
# マイグレーション関連ターゲット
# ================================================

# migrationファイルを生成する。
# e.g.
# make new DB=garbage_category_rule_quiz NAME=users
# make new DB=garbage_category_rule_quiz NAME=add_column_xxxx_to_users
new: .check-db .check-name
	migrate create -ext sql -dir db/$(DB)/migrations $(NAME)

# マイグレーションを適用する。
# e.g.
# make up DB=garbage_category_rule_quiz
up: .check-db
	MODE=up TARGET_DB=$(DB) go run main.go

# マイグレーションを一世代前にロールバックする
# e.g.
# make down DB=garbage_category_rule_quiz
down: .check-db
	MODE=down TARGET_DB=$(DB) go run main.go

# ダンプファイルを作成する。
# e.g.
# make dump DB=garbage_category_rule_quiz
dump: .check-db
	docker compose exec db mysqldump -u $$(docker compose exec -T db printenv MYSQL_USER | tr -d '\r') -p$$(docker compose exec -T db printenv MYSQL_PASSWORD | tr -d '\r') --databases $(DB) > db/$(DB)/dump.sql

# dumpファイルからseedデータを復元する。
# e.g.
# make restore DB=garbage_category_rule_quiz
restore: .check-db
	docker compose exec -T db mysql -u $$(docker compose exec -T db printenv MYSQL_USER | tr -d '\r') -p$$(docker compose exec -T db printenv MYSQL_PASSWORD | tr -d '\r') -e "DROP DATABASE IF EXISTS $(DB); CREATE DATABASE $(DB) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
	docker compose exec -T db mysql -u $$(docker compose exec -T db printenv MYSQL_USER | tr -d '\r') -p$$(docker compose exec -T db printenv MYSQL_PASSWORD | tr -d '\r') $(DB) < db/$(DB)/dump.sql

# ローカル環境でマイグレーションを実行する。
# e.g.
# make run-local-db DB=garbage_category_rule_quiz
run-local-db:
	docker compose up -d
	sleep 5
	go run main.go TARGET_DB=$(DB)

# ================================================
# ECR デプロイ用ターゲット
# ================================================

# ローカルでDockerイメージをビルドする
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

# ECSタスクをスタンドアロンで実行する（マイグレーション実行）
# e.g. make run-ecs ENV=stg DB=garbage_category_rule_quiz
run-ecs: .check-db
	aws ecs run-task \
		--cluster $(CLUSTER) \
		--task-definition $(SERVICE_NAME)-$(ENV) \
		--launch-type FARGATE \
		--network-configuration "awsvpcConfiguration={subnets=[$(SUBNET_ID)],securityGroups=[$(SECURITY_GROUP_ID)],assignPublicIp=DISABLED}" \
		--overrides "{\"containerOverrides\":[{\"name\":\"app\",\"environment\":[{\"name\":\"MODE\",\"value\":\"up\"},{\"name\":\"TARGET_DB\",\"value\":\"$(DB)\"}]}]}" \
		--region $(AWS_REGION)

# ローカルで実行する（Docker経由）
# e.g. make run-local
run-local:
	docker run --rm \
		-e DB_HOST=$(DB_HOST) \
		-e DB_PORT=$(DB_PORT) \
		-e DB_USER=$(DB_USER) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e TARGET_DB=$(TARGET_DB) \
		-e MODE=$(MODE) \
		$(SERVICE_NAME):latest

# クリーンアップ
# e.g. make clean
clean:
	docker rmi $(SERVICE_NAME):$(GIT_COMMIT_HASH) $(SERVICE_NAME):latest || true
	docker rmi $(IMAGE_REPOSITORY):$(GIT_COMMIT_HASH) $(IMAGE_REPOSITORY):latest || true

.PHONY: new up down dump restore run-local-db build ecr-login push deploy run-ecs run-local clean

.check-name:
	@if [ -z "$(NAME)" ]; then \
		echo "NAME is required"; \
		exit 1; \
	fi

.check-db:
	@if [ -z "$(DB)" ]; then \
		echo "DB_NAME is required"; \
		exit 1; \
	fi
