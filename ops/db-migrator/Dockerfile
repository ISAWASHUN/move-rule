FROM golang:1.25.5 AS builder
WORKDIR /app

# 依存関係キャッシュを効かせるために先に go.mod だけコピー
COPY go.mod go.sum ./
RUN go mod download

# アプリコードをコピー
COPY . .

# 本番用バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o db-migrator

FROM gcr.io/distroless/base-debian12
WORKDIR /app

# distrolessでは65532が非rootユーザーとして推奨されている
USER 65532

# builderからビルド成果物だけコピー
COPY --from=builder /app/db-migrator .

# メインコマンド（固定）
ENTRYPOINT ["./db-migrator"]
