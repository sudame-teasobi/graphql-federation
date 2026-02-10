# GraphQL Federation Teasobi

GraphQL Federation を学習・実験するためのリポジトリ。

## プロジェクト構成

2つの独立した GraphQL サービスで構成されています。

| サービス | ポート | 説明 |
|---------|--------|------|
| User Service | 8080 | ユーザー情報の管理 |
| Task Service | 8081 | タスク情報の管理 |

## 技術スタック

- Go 1.25
- [gqlgen](https://github.com/99designs/gqlgen) v0.17.86 - GraphQL サーバー生成
- [mise](https://mise.jdx.dev/) - ツールバージョン管理

## ディレクトリ構造

```
.
├── cmd/
│   ├── user/
│   │   ├── main.go            # User サービスのエントリーポイント
│   │   └── users.json         # 初期ユーザーデータ
│   └── task/
│       ├── main.go            # Task サービスのエントリーポイント
│       └── tasks.json         # 初期タスクデータ
├── internal/
│   ├── user/
│   │   ├── schema.graphql     # User サービスの GraphQL スキーマ
│   │   ├── repository.go      # データリポジトリ
│   │   ├── id.go              # タイプ定数
│   │   └── graph/
│   │       ├── generated.go   # 生成コード
│   │       ├── model/
│   │       │   └── models_gen.go
│   │       └── resolver/
│   │           ├── resolver.go
│   │           └── schema.resolvers.go
│   └── task/
│       ├── schema.graphql     # Task サービスの GraphQL スキーマ
│       ├── repository.go      # データリポジトリ
│       ├── id.go              # タイプ定数
│       └── graph/
│           ├── generated.go   # 生成コード
│           ├── model/
│           │   └── models_gen.go
│           └── resolver/
│               ├── resolver.go
│               └── schema.resolvers.go
├── gqlgen.user.yml            # User サービスの gqlgen 設定
├── gqlgen.todo.yml            # Task サービスの gqlgen 設定
├── generate.go                # コード生成エントリーポイント
└── mise.toml                  # mise 設定
```

## セットアップ

```bash
# ツールのインストール（mise 使用時）
mise install

# 依存関係のインストール
go mod download
```

## 使い方

### コード生成

GraphQL スキーマを変更した後、以下のコマンドでコードを再生成します。

```bash
go generate
```

### サービスの起動

```bash
# User Service（別ターミナルで実行）
go run ./cmd/user

# Task Service（別ターミナルで実行）
go run ./cmd/task
```

### GraphQL Playground

サービス起動後、ブラウザでアクセスできます。

- User Service: http://localhost:8080/
- Task Service: http://localhost:8081/

### クエリ例

**User Service:**
```graphql
query {
  users {
    id
    name
  }
}

query {
  user(id: "user:0") {
    id
    name
  }
}
```

**Task Service:**
```graphql
query {
  tasks {
    id
    title
  }
}

query {
  task(id: "task:0") {
    id
    title
  }
}
```

## GraphQL スキーマ

### User Service

```graphql
interface Node {
  id: ID!
}

type User implements Node {
  id: ID!
  name: String!
}

type Query {
  users: [User!]!
  user(id: ID!): User
  node(id: ID!): Node
}
```

### Task Service

```graphql
interface Node {
  id: ID!
}

type Task implements Node {
  id: ID!
  title: String!
}

type Query {
  node(id: ID!): Node
  task(id: ID!): Task
  tasks: [Task!]!
}
```

## 今後の予定

- [ ] GraphQL Federation の有効化
- [ ] Apollo Router / Gateway の導入
- [ ] サービス間のエンティティ参照（Task に User フィールドを追加）
