# GraphQL Federation Teasobi

GraphQL Federation を学習・実験するためのリポジトリ。Go + gqlgen で構成された2つの GraphQL サブグラフサービスを Apollo Federation v2 で統合している。

## プロジェクト構成

2つのサブグラフサービスと、それらを統合する Apollo Router で構成されています。

| サービス | ポート | 説明 |
|---------|--------|------|
| User Service | 8080 | ユーザー情報の管理（`User` エンティティを所有） |
| Task Service | 8081 | タスク情報の管理（`Task` エンティティを所有） |
| Apollo Router | 4000 | サブグラフを統合するゲートウェイ |

### エンティティ関係

- **User Service**: `User` エンティティを所有。`User.tasks` フィールドで Task を参照
- **Task Service**: `Task` エンティティを所有。`Task.user` フィールドで User を参照

各サービスは他サービスのエンティティをスタブ型（`resolvable: false`）として定義し、`@key(fields: "id")` で Federation のエンティティ解決をサポートしている。

## 技術スタック

- Go 1.25
- [gqlgen](https://github.com/99designs/gqlgen) v0.17.86 - GraphQL サーバー生成（Federation v2 対応）
- [Apollo Router](https://www.apollographql.com/docs/router/) - Federation ゲートウェイ
- [Rover](https://www.apollographql.com/docs/rover/) - スーパーグラフのコンポジション
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
│   │   ├── schema.graphql     # User サブグラフの GraphQL スキーマ
│   │   ├── repository.go      # データリポジトリ
│   │   ├── id.go              # タイプ定数
│   │   └── graph/
│   │       ├── generated.go   # 自動生成コード（編集禁止）
│   │       ├── federation.go  # 自動生成 Federation ランタイム（編集禁止）
│   │       ├── model/
│   │       │   └── models_gen.go  # 自動生成モデル（編集禁止）
│   │       └── resolver/
│   │           ├── resolver.go
│   │           ├── schema.resolvers.go
│   │           └── entity.resolvers.go  # Federation エンティティリゾルバ
│   └── task/
│       ├── schema.graphql     # Task サブグラフの GraphQL スキーマ
│       ├── repository.go      # データリポジトリ
│       ├── id.go              # タイプ定数
│       └── graph/
│           ├── generated.go   # 自動生成コード（編集禁止）
│           ├── federation.go  # 自動生成 Federation ランタイム（編集禁止）
│           ├── model/
│           │   └── models_gen.go  # 自動生成モデル（編集禁止）
│           └── resolver/
│               ├── resolver.go
│               ├── schema.resolvers.go
│               └── entity.resolvers.go  # Federation エンティティリゾルバ
├── supergraph.yaml            # Apollo Router のサブグラフ構成
├── supergraph.graphql         # コンポーズ済みスーパーグラフスキーマ
├── gqlgen.user.yml            # User サービスの gqlgen 設定
├── gqlgen.task.yml            # Task サービスの gqlgen 設定
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

GraphQL スキーマを変更した後、以下のコマンドでコードを再生成します。gqlgen によるコード生成に加え、Rover によるスーパーグラフのコンポジションも実行されます。

```bash
go generate
```

### サービスの起動

```bash
# User Service（別ターミナルで実行）
go run ./cmd/user

# Task Service（別ターミナルで実行）
go run ./cmd/task

# Apollo Router（別ターミナルで実行、両サブグラフ起動後）
router --supergraph supergraph.graphql
```

### GraphQL Playground

サービス起動後、ブラウザでアクセスできます。

- User Service: http://localhost:8080/
- Task Service: http://localhost:8081/

### クエリ例

**User Service（直接）:**
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

**Task Service（直接）:**
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

**Apollo Router 経由（Federation クエリ）:**
```graphql
# ユーザーとそのタスクを一括取得
query {
  users {
    id
    name
    tasks {
      id
      title
    }
  }
}

# タスクと担当ユーザーを一括取得
query {
  tasks {
    id
    title
    user {
      id
      name
    }
  }
}
```

## GraphQL スキーマ

### User Service（サブグラフ）

```graphql
extend schema @link(url: "https://specs.apollo.dev/federation/v2.3", import: ["@key", "@shareable"])

interface Node {
  id: ID!
}

type User implements Node @key(fields: "id") {
  id: ID!
  name: String!
  tasks: [Task!]!
}

type Task @key(fields: "id", resolvable: false) {
  id: ID!
}

type Query {
  users: [User!]!
  user(id: ID!): User
}
```

### Task Service（サブグラフ）

```graphql
extend schema @link(url: "https://specs.apollo.dev/federation/v2.3", import: ["@key", "@shareable"])

interface Node {
  id: ID!
}

type Task implements Node @key(fields: "id") {
  id: ID!
  title: String!
  user: User!
}

type User @key(fields: "id", resolvable: false) {
  id: ID!
}

type Query {
  task(id: ID!): Task
  tasks: [Task!]!
}
```
