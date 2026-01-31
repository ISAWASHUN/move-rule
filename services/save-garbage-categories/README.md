# save-garbage-categories

`fetch-garbage-categories`で取得したJSONファイルをMySQLデータベースに保存するサービスです。

## 機能

- JSONファイルからゴミ品目データを読み込み
- 自治体（municipalities）テーブルへの保存
- 分別区分（garbage_categories）テーブルへの保存
- ゴミ品目（garbage_items）テーブルへの保存

## ディレクトリ構成

```
save-garbage-categories/
├── cmd/
│   └── main.go                 # エントリーポイント
├── internal/
│   ├── domain/
│   │   ├── garbage_item.go     # ドメインモデル
│   │   └── repository.go       # リポジトリインターフェース
│   ├── infrastructure/
│   │   ├── repository/
│   │   │   ├── entity.go           # GORMエンティティ
│   │   │   ├── municipality.go     # 自治体リポジトリ
│   │   │   ├── garbage_category.go # 分別区分リポジトリ
│   │   │   └── garbage_item.go     # ゴミ品目リポジトリ
│   │   └── storage/
│   │       └── file.go         # JSONファイル読み込み
│   └── usecase/
│       └── save.go             # 保存ユースケース
├── go.mod
└── README.md
```

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `INPUT_FILE` | `../fetch-garbage-categories/internal/infrastructure/storage/file/latest.json` | 入力JSONファイルのパス（ローカル） |
| `S3_BUCKET` | - | S3バケット名（指定するとS3から読み込み） |
| `S3_KEY` | `data/latest.json` | S3オブジェクトキー |
| `S3_PREFIX` | `data` | S3プレフィックス（`S3_KEY`未指定時に使用） |
| `AWS_REGION` | `ap-northeast-1` | AWSリージョン |
| `DB_HOST` | `localhost` | MySQLホスト |
| `DB_PORT` | `3306` | MySQLポート |
| `DB_USER` | `root` | MySQLユーザー |
| `DB_PASSWORD` | `password` | MySQLパスワード |
| `DB_NAME` | `garbage_category_rule_quiz` | データベース名 |
| `LOG_LEVEL` | `info` | ログレベル (debug, info, warn, error) |

## 使用方法

### ビルド

```bash
cd services/save-garbage-categories
go build -o save-garbage-categories cmd/main.go
```

### 実行

```bash
# デフォルト設定で実行（fetch-garbage-categoriesのlatest.jsonを使用）
./save-garbage-categories

# カスタムファイルパスを指定
INPUT_FILE=/path/to/data.json ./save-garbage-categories

# S3から読み込み
S3_BUCKET=garbage-category-rule-data-stg S3_PREFIX=data ./save-garbage-categories
S3_BUCKET=garbage-category-rule-data-stg S3_KEY=data/latest.json ./save-garbage-categories

# DB接続設定を指定
DB_HOST=localhost DB_PORT=3306 DB_USER=root DB_PASSWORD=password ./save-garbage-categories
```

## データフロー

1. `fetch-garbage-categories` で外部APIからデータを取得し、JSONファイルに保存
2. `save-garbage-categories` でJSONファイルを読み込み、DBに保存
3. `quiz` サービスでDBからデータを取得してクイズを提供

```
[外部API] → [fetch-garbage-categories] → [JSON] → [save-garbage-categories] → [MySQL]
                                                                                  ↓
                                                            [quiz] ← ─ ─ ─ ─ ─ ─ ┘
```
