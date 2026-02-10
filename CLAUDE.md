# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

GraphQL Federation を学習・実験するためのリポジトリ。Go + gqlgen で構成された2つの独立した GraphQL サブグラフサービス（User / Task）を持つ。Federation はまだ有効化されていない（gqlgen.*.yml 内の federation セクションがコメントアウト状態）。

## よく使うコマンド

```bash
# コード生成（GraphQL スキーマ変更後に実行）
go generate

# サービス起動
go run ./cmd/user   # User Service (port 8080)
go run ./cmd/task   # Task Service (port 8081)

# ツールインストール（mise）
mise install
```

## アーキテクチャ

- **Go モジュール名**: `gft`
- **gqlgen 設定**: サービスごとに個別の設定ファイル（`gqlgen.user.yml`, `gqlgen.task.yml`）
- **コード生成エントリーポイント**: `generate.go`（`//go:generate` ディレクティブで両サービスの gqlgen を実行）

### サービス構成

各サービスは同一パターンで構成されている:

```
cmd/<service>/main.go          # エントリーポイント、embed で JSON データを読み込み
cmd/<service>/<data>.json      # 初期データ（go:embed で埋め込み）
internal/<service>/
  schema.graphql               # GraphQL スキーマ定義（手動編集）
  id.go                        # Type 定数（Node ID のプレフィックス）
  repository.go                # インメモリデータリポジトリ
  model/models_gen.go          # gqlgen 自動生成モデル（編集禁止）
  graph/generated.go           # gqlgen 自動生成コード（編集禁止）
  graph/resolver/resolver.go   # DI 用 Resolver 構造体（手動編集）
  graph/resolver/schema.resolvers.go  # リゾルバ実装（手動編集、生成時にマージされる）
```

### ID 体系

Node インターフェースの ID は `<type>:<index>` 形式（例: `user:0`, `task:0`）。`strings.Cut(id, ":")` でタイプを抽出し、`node` クエリのディスパッチに使用。

### 手動編集ファイル vs 自動生成ファイル

- **手動編集**: `schema.graphql`, `resolver.go`, `schema.resolvers.go`, `repository.go`, `id.go`, `main.go`
- **自動生成（編集禁止）**: `generated.go`, `models_gen.go`
