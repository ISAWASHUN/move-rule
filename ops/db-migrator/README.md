# db-migrator

データベースマイグレーションを管理・実行するためのツールです。

## 前提条件

- Docker
- Docker Compose
- golang-migrate (ローカル環境用)

## セットアップ

### ローカル環境への golang-migrate のインストール

マイグレーションファイルの生成はローカル環境で行います。以下のコマンドで golang-migrateをインストールしてください。

```bash
brew install golang-migrate
```

### mysql コンテナを起動

mysqlのコンテナをデーモン起動します。

```bash
docker compose up -d
```

### db-migrator を実行

db-migrator を実行します。

```bash
go run main.go

>
2026/01/17 20:55:47 INFO start db-migrator
2026/01/17 20:55:47 INFO mode !BADKEY=up
2026/01/17 20:55:47 INFO dbName !BADKEY=garbage_category_rule_quiz
2026/01/17 20:55:47 INFO connected to default db
2026/01/17 20:55:47 INFO Migrations applied successfully
```

※ `TARGET_DB` を指定することで、指定した DB に対してマイグレーションを実行できます。指定しない場合は全ての DB に対してマイグレーションを実行します。
※ `MODE` を `up` にするとマイグレーションを実行します。`down` にするとマイグレーションをロールバックします。

## 運用

### マイグレーションファイルの作成

新しいマイグレーションファイルを生成する場合は、以下のコマンドを実行します。

```bash
make new DB=[DB名] NAME=[マイグレーションファイル名]
```

例：

```bash
make new DB=garbage_category_rule_quiz NAME=create_users
```

これにより、`db/[DB名]/migrations/`ディレクトリに以下の 2つのファイルが生成されます。

- `YYYYMMDDHHMMSS_create_users.up.sql`
- `YYYYMMDDHHMMSS_create_users.down.sql`

マイグレーションファイル内には SQL 文を記述してください。

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
```

### マイグレーションファイルの実行(最新バージョンまで)

マイグレーションファイルを実行する場合は、以下のコマンドを実行します.

```bash
make up DB=[DB名]
```

例：

```bash
make up DB=garbage_category_rule_quiz
```

### マイグレーションのロールバック(1 つのバージョンのみ)

マイグレーションファイルをロールバックする場合は、以下のコマンドを実行します。

```bash
make down DB=[DB名]
```

## 環境変数について

環境変数で以下の項目を設定できます。

- `DB_HOST` : DB インスタンスのホスト名
- `DB_PORT` : DB インスタンスのポート番号
- `DB_USER` : migration を実行するユーザー名
- `DB_PASSWORD` : migration を実行するユーザのパスワード
- `TARGET_DB`: マイグレーションを実行するDB名を指定します。指定しない場合は全ての DB に対して実行します。
