## 要件

ユーザーがゴミを選択して、分類方法の差分を出力するアプリケーション

## データ要件

- 東京都のデータセットを使用
  - 【板橋区】ゴミ分別のデータです。
    - https://service.api.metro.tokyo.lg.jp/api/t131199d3000000001-10af70080e2503877feb2bf2c9a42171-0/json
  - 【立川市】ごみの分別方法一覧
    - https://service.api.metro.tokyo.lg.jp/api/t132021d3000000001-ef24963d14f44ffddea8de17cb2d0ea6-0/json

## 画面要件

- 入力画面
  - ユーザーが自治体を入力する
  - ユーザーがセレクトボックスからゴミを選択する
- 検索結果画面
  - ユーザーが検索した結果を返す画面

## 設計

### データ取得

1. ECSのスケジュールされたタスクでcron形式で実行する

### データ処理

- 入力値

板橋、立川市、ゴミの名前

- 出力値

ゴミの分類方法

### 技術選定

### バックエンド

- Go
  - Gin
  - slog
  - go-migrator
  - gorm or sqlc
  - golangci-lint
  - swag

### フロントエンド

- React

### インフラ

- Docker
  - MySQL (ローカル)
  - Go (Lambda用)
- CI/CD
    - GitHub Actions
- AWS
  - Amplify
  - Lambda
    - サーバーの置き場所
  - RDS
    - MySQL

## DB設計
```mermaid
erDiagram
    municipalities ||--o{ garbage_items : "belongs_to"
    waste_categories ||--o{ garbage_items : "categorizes"

    municipalities {
        int id PK "AUTO_INCREMENT"
        int code UK "全国地方公共団体コード"
        varchar name "地方公共団体名"
        datetime created_at "レコード作成日時"
        datetime updated_at "レコード更新日時"
    }

    waste_categories {
        int id PK "AUTO_INCREMENT"
        varchar name UK "分別区分名"
        datetime created_at "レコード作成日時"
        datetime updated_at "レコード更新日時"
    }

    garbage_items {
        int id PK "AUTO_INCREMENT"
        int municipality_id FK "地方公共団体ID"
        int waste_category_id FK "分別区分ID"
        varchar area_name "地区名"
        varchar item_name "ゴミの品目"
        varchar item_name_kana "ゴミの品目_カナ"
        varchar item_name_english "ゴミの品目_英字"
        text notes "注意点"
        text remarks "備考"
        int bulk_garbage_fee "粗大ごみ回収料金"
        datetime created_at "レコード作成日時"
        datetime updated_at "レコード更新日時"
    }
```
