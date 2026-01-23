# Swagger ドキュメント生成

このディレクトリには、swaggoによって生成されるSwaggerドキュメントが配置されます。

## セットアップ

### 1. swag コマンドのインストール

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Swagger ドキュメントの生成

プロジェクトルート（`services/quiz`）で以下のコマンドを実行:

```bash
swag init -g cmd/main.go -o docs
```

### 3. サーバーの起動

```bash
go run cmd/main.go
```

### 4. Swagger UI へのアクセス

ブラウザで以下のURLにアクセス:

```
http://localhost:8080/swagger/index.html
```

## 注意事項

- `swag init` を実行すると、このディレクトリに以下のファイルが生成されます:
  - `docs.go` - Swagger仕様のGoコード
  - `swagger.json` - JSON形式のSwagger仕様
  - `swagger.yaml` - YAML形式のSwagger仕様

- コードを変更した後は、再度 `swag init` を実行してドキュメントを更新してください。
