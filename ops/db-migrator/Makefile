ROOT_DIR = ../../

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
# make run-local
run-local:
	docker compose up -d
	sleep 5
	go run main.go TARGET_DB=$(DB)

PHONY: new up down dump restore run-local

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
