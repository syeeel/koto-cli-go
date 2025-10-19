# Claude Code Go開発テンプレート

Claude CodeのDevContainer環境でGo言語開発を行うためのテンプレートです。
DevContainerを使用することで、ローカル環境を汚さずに、すぐに開発を始められます。

## 特徴

- **Go 1.25.3** - 最新の安定版Go
- **Claude Code完全統合** - AIアシスタントによる開発支援
- **開発ツール完備**
  - `gopls` - Go Language Server（コード補完、リファクタリング、型チェック）
  - `golangci-lint` - 統合リンター（複数のリンターを統合）
  - `delve (dlv)` - デバッガー
  - `air` - ホットリロード開発サーバー
- **VS Code / Cursor最適化** - Go開発に必要な拡張機能とエディタ設定を事前構成
- **開発環境の完全分離** - DevContainerによる再現可能な開発環境

## セットアップ

### 前提条件

- Docker Desktop がインストールされていること
- VS Code または Cursor エディタがインストールされていること
- VS Code / Cursorに「Dev Containers」拡張機能がインストールされていること

### 1. リポジトリのクローン

```bash
git clone <このリポジトリのURL>
cd <リポジトリ名>
```

### 2. DevContainerでの起動

VS Code / Cursorでこのプロジェクトを開き、以下のいずれかの方法でDevContainerを起動：

**方法1: 通知から起動**
- エディタ右下に表示される「Reopen in Container」をクリック

**方法2: コマンドパレットから起動**
1. `Cmd+Shift+P` (Mac) / `Ctrl+Shift+P` (Windows/Linux) でコマンドパレットを開く
2. `Dev Containers: Rebuild and Reopen in Container` を選択

初回起動時はDockerイメージのビルドに数分かかります。

### 3. 動作確認

コンテナが起動したら、ターミナルで以下のコマンドを実行して環境を確認：

```bash
# Goのバージョン確認
go version

# 開発ツールの確認
gopls version
golangci-lint --version
dlv version
air -v
```

## 使い方

### 基本的な開発

```bash
# アプリケーションの実行
go run main.go

# ビルド
go build -o app

# テストの実行
go test ./...

# リンターの実行
golangci-lint run
```

### ホットリロード（Air使用）

```bash
# Airでホットリロードを有効にして開発
air

# .air.toml で設定をカスタマイズ可能
```

ファイルを保存すると自動的に再ビルド・再起動されます。

### デバッグ

VS Code / Cursorのデバッグ機能を使用するか、コマンドラインから：

```bash
# デバッガーを起動
dlv debug
```

## プロジェクト構成

```
.
├── .devcontainer/          # DevContainer設定ファイル
│   ├── Dockerfile          # Go開発環境のDockerイメージ定義
│   └── devcontainer.json   # DevContainerの設定
├── .claude/                # Claude Code設定
├── .vscode/                # VS Code設定
├── .air.toml              # Airホットリロード設定
├── main.go                # サンプルHTTPサーバー
├── go.mod                 # Go module定義
├── README.md              # このファイル
└── .gitignore             # Gitignore設定
```

## サンプルアプリケーション

このテンプレートには簡単なHTTPサーバーが含まれています：

```bash
# サーバーを起動
go run main.go

# 別のターミナルで確認
curl http://localhost:8080
curl http://localhost:8080/health
```

## 開発環境の詳細

### インストール済みツール

- **Go 1.23.2** - 最新の安定版Go
- **Node.js 20** - Claude Code実行環境
- **Git & GitHub CLI** - バージョン管理
- **Zsh** - デフォルトシェル（高機能なコマンドライン）
- **Git Delta** - 美しい差分表示

### Go開発ツール

- **gopls** - Go Language Server Protocol実装（IntelliSense、リファクタリング、型チェック）
- **golangci-lint** - 統合リンター（go vet、staticcheck等、複数のリンターを統合）
- **delve (dlv)** - Goデバッガー（ブレークポイント、ステップ実行）
- **air** - ライブリロードツール（ファイル変更を検知して自動再起動）

### VS Code / Cursor拡張機能

自動的にインストールされる拡張機能：

- **Claude Code** - AIアシスタント
- **Go** - Go言語サポート
- **GitLens** - Git強化ツール

### エディタ設定

`.devcontainer/devcontainer.json`で以下が自動設定されます：

- 保存時の自動フォーマット（`gofmt`）
- 保存時のインポート自動整理
- golangci-lintによる自動リント
- Zshをデフォルトターミナルに設定

## カスタマイズ

### Goバージョンの変更

`.devcontainer/Dockerfile` の `FROM golang:1.23.2-bookworm` を変更してください。

### 追加のGo依存関係

```bash
go get <package-name>
```

### 追加の開発ツール

`.devcontainer/Dockerfile` の該当セクションに追加してください。

## トラブルシューティング

### コンテナが起動しない

```bash
# コンテナを完全に再ビルド
Dev Containers: Rebuild Container Without Cache
```

### Go言語サーバーが動作しない

```bash
# goplsの再インストール
go install golang.org/x/tools/gopls@latest
```

## Claude Code について

このテンプレートは**Claude Code**との統合を前提に設計されています。

### Claude Codeとは？

Claude CodeはAnthropicが開発したAIアシスタントで、VS Code / Cursor内で直接動作し、コードの作成、リファクタリング、バグ修正、説明などを支援します。

### このテンプレートでできること

- **コード生成**: 自然言語でリクエストして、Goコードを生成
- **バグ修正**: エラーメッセージを共有して、修正案を取得
- **コードレビュー**: 既存コードの改善提案を受ける
- **テスト作成**: 関数のユニットテストを自動生成
- **リファクタリング**: コードの改善提案と実装

### 使い方

1. DevContainerでプロジェクトを開く
2. Claude Codeを起動（エディタ内のClaude Code拡張機能から）
3. 自然言語でリクエストを入力
4. Claude Codeがコードを生成・編集

詳しくは[Claude Code公式ドキュメント](https://docs.anthropic.com/en/docs/claude-code/overview)をご覧ください。

## よくある質問

### Q: ローカル環境にGoをインストールする必要はありますか？

A: いいえ、必要ありません。全ての開発環境はDevContainer内に含まれているため、Dockerさえあれば開発を始められます。

### Q: このテンプレートは商用プロジェクトで使用できますか？

A: はい、自由に使用できます。詳細はLICENSE.mdをご確認ください。

### Q: Go以外の言語でも同様のテンプレートはありますか？

A: Claude Codeは様々な言語に対応しています。詳しくは公式ドキュメントをご覧ください。

### Q: DevContainerを使わずに使用できますか？

A: このテンプレートはDevContainer前提で設計されていますが、Dockerfileを参考にローカル環境を構築することも可能です。

## 参考リンク

- [Go公式ドキュメント](https://go.dev/doc/) - Go言語の公式ドキュメント
- [Claude Code公式サイト](https://claude.ai/code) - Claude Codeの詳細
- [Claude Code ドキュメント](https://docs.anthropic.com/en/docs/claude-code/overview) - 使い方ガイド
- [golangci-lint](https://golangci-lint.run/) - Goリンターツール
- [Air](https://github.com/air-verse/air) - ホットリロードツール
- [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers) - VS Code Dev Containers

## ライセンス

このプロジェクトのライセンスについては、[LICENSE.md](LICENSE.md)をご確認ください。

## セキュリティ

セキュリティに関する問題を発見した場合は、[SECURITY.md](SECURITY.md)をご確認ください。
