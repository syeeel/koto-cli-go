# koto - ToDo管理CLI

**koto**（琴）は、Go言語で開発されたインタラクティブなToDo管理CLIツールです。
[Bubbletea](https://github.com/charmbracelet/bubbletea)フレームワークを使用した美しいターミナルUIで、快適なタスク管理体験を提供します。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)

## ✨ 特徴

- 🎨 **リッチなTUI** - Bubbletea/Lipglossによる美しいターミナルインターフェース
- ⚡ **軽量・高速** - Pure Go（CGO不要）で高速起動
- 📊 **優先度管理** - 3段階の優先度設定（🔴高 🟡中 🟢低）
- 📅 **期限管理** - 期限日の設定と期限切れ警告
- 💾 **SQLite保存** - ローカルデータベースで確実にデータを保持
- 📤 **エクスポート/インポート** - JSON形式でバックアップ・移行可能
- ⌨️ **Vimライクなキーバインド** - j/kでの快適なナビゲーション
- 🔍 **ステータスフィルター** - 未完了/完了済みで絞り込み表示
- 🍅 **ポモドーロタイマー** - 25分間のタイマーで集中作業をサポート、作業時間の自動記録

## 📦 インストール

### インストールスクリプト（推奨・macOS/Linux）

最も簡単な方法：

```bash
curl -sSfL https://raw.githubusercontent.com/syeeel/koto-cli-go/main/install.sh | sh
```

このスクリプトは：
- 最新バージョンを自動検出
- お使いのOS/アーキテクチャに合わせてダウンロード
- `~/.local/bin`にインストール
- PATHの設定方法を案内

### Homebrew（macOS - 準備中🚧）

**現在セットアップ中です。次回リリース（v1.0.1以降）から利用可能になります。**

準備が完了次第、以下のコマンドでインストールできるようになります：

```bash
brew tap syeeel/tap
brew install koto
```

**それまでは、インストールスクリプト（上記）をご利用ください。**

### Go install

Go環境がある場合：

```bash
go install github.com/syeeel/koto-cli-go/cmd/koto@latest
```

### ビルド済みバイナリ

[Releases](https://github.com/syeeel/koto-cli-go/releases/latest)ページから、お使いのプラットフォーム向けのバイナリをダウンロードできます。

対応プラットフォーム：
- **macOS**: darwin_amd64 (Intel), darwin_arm64 (Apple Silicon)
- **Linux**: linux_amd64, linux_arm64
- **Windows**: windows_amd64

ダウンロード後、解凍してPATHの通った場所に配置してください。

### ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/syeeel/koto-cli-go.git
cd koto-cli-go

# 依存関係をダウンロード
go mod download

# ビルド
go build -o bin/koto ./cmd/koto

# 実行
./bin/koto
```

## 🚀 使い方

### アプリケーションの起動

```bash
koto
```

起動すると、インタラクティブなTUIが表示されます。

### 基本コマンド

#### ToDoの追加

```bash
/add 買い物に行く
/add レポート作成 --desc="第5章をまとめる" --priority=high --due=2025-10-25
```

**オプション**:
- `--desc="説明文"` - ToDoの詳細説明
- `--priority=low|medium|high` - 優先度（デフォルト: medium）
- `--due=YYYY-MM-DD` - 期限日

#### ToDoの表示

```bash
/list                      # すべてのToDoを表示
/list --status=pending     # 未完了のみ
/list --status=completed   # 完了済みのみ
```

#### ToDoの完了

```bash
/done 1    # ID 1のToDoを完了にする
```

#### ToDoの編集

```bash
/edit 1 --title="新しいタイトル"
/edit 1 --priority=high
/edit 1 --desc="新しい説明"
/edit 1 --due=2025-12-31
```

#### ToDoの削除

```bash
/delete 1    # ID 1のToDoを削除
```

#### エクスポート/インポート

```bash
/export ~/my-todos.json     # JSONファイルにエクスポート
/import ~/todos-backup.json # JSONファイルからインポート
```

#### ヘルプ

```bash
/help    # ヘルプ画面を表示
```

#### ポモドーロタイマー

```bash
/pomo              # 25分間のタイマーを開始
/pomo 1            # ID 1のToDoに紐づけて25分間のタイマーを開始（作業時間を自動記録）
```

**ポモドーロタイマーの使い方**:
- タイマー実行中は専用画面が表示されます
- 25分経過するとアラーム音が鳴ります
- タスクIDを指定した場合、作業時間が自動的に記録されます
- `Esc`キーでタイマーをキャンセルしてメイン画面に戻ります

### ⌨️ キーボードショートカット

| キー | 動作 |
|------|------|
| `↑` / `k` | カーソルを上に移動 |
| `↓` / `j` | カーソルを下に移動 |
| `Enter` | コマンドを実行 |
| `Esc` | 入力欄をクリア |
| `?` | ヘルプ画面を表示/非表示 |
| `Ctrl+C` | アプリケーション終了 |

### 📺 画面の見方

```
📝 koto - ToDo Manager

  ⬜ 🔴 [1] 重要な会議の準備
> ✅ 🟡 [2] 買い物リスト ⚠ OVERDUE
  ⬜ 🟢 [3] メールの返信

> /add 新しいタスク

Commands: /add, /list, /done, /delete, /edit, /help | Navigate: ↑/↓ or j/k | Help: ? | Quit: Ctrl+C
```

**表示の説明**:
- `>` - 現在選択中のToDo（カーソル位置）
- `⬜` - 未完了
- `✅` - 完了済み
- `🔴🟡🟢` - 優先度（高・中・低）
- `[数字]` - ToDo ID
- `⚠ OVERDUE` - 期限切れの警告
- `🍅 XXXm` - 累積作業時間（ポモドーロタイマーで記録）

## 📁 データの保存場所

すべてのToDoは以下のSQLiteデータベースに保存されます：

```
~/.koto/koto.db
```

バックアップを取る場合は、このファイルをコピーするか、`/export` コマンドを使用してください。

## 🏗️ アーキテクチャ

kotoは、クリーンアーキテクチャに基づいたレイヤー構造を採用しています：

```
┌─────────────────┐
│   TUI Layer     │  Bubbletea UI（コマンド入力、表示）
├─────────────────┤
│ Service Layer   │  ビジネスロジック、バリデーション
├─────────────────┤
│Repository Layer │  データアクセス（SQLite）
├─────────────────┤
│  Model Layer    │  データ構造定義
└─────────────────┘
```

### ディレクトリ構成

```
koto-cli-go/
├── cmd/
│   └── koto/              # メインエントリーポイント
│       └── main.go
├── internal/
│   ├── model/             # データモデル（Todo, Status, Priority）
│   │   ├── todo.go
│   │   └── todo_test.go
│   ├── repository/        # データアクセス層
│   │   ├── repository.go
│   │   ├── sqlite.go
│   │   └── sqlite_test.go
│   ├── service/           # ビジネスロジック層
│   │   ├── todo_service.go
│   │   └── todo_service_test.go
│   ├── tui/               # ターミナルUI
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── views.go
│   │   ├── styles.go
│   │   └── commands.go
│   └── config/            # 設定管理
│       └── config.go
├── migrations/            # データベーススキーマ
│   └── 001_init.sql
├── docs/                  # ドキュメント
│   ├── design/            # 設計書
│   └── implementation/    # 実装管理
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ 開発環境

### 必要な環境

- Go 1.21以上
- SQLite 3（Pure Go実装を使用するため不要）

### 依存関係

- [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUIフレームワーク
- [github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) - TUIコンポーネント
- [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - スタイリング
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite

### 開発コマンド

```bash
# 依存関係のダウンロード
go mod download

# テストの実行
go test ./...
go test -v ./internal/model/...     # Model層のみ
go test -v ./internal/repository/... # Repository層のみ
go test -v ./internal/service/...    # Service層のみ

# リント
go vet ./...
golangci-lint run  # golangci-lintがインストール済みの場合

# ビルド
go build -o bin/koto ./cmd/koto

# クロスコンパイル
GOOS=darwin GOARCH=amd64 go build -o bin/koto-darwin-amd64 ./cmd/koto
GOOS=linux GOARCH=amd64 go build -o bin/koto-linux-amd64 ./cmd/koto
GOOS=windows GOARCH=amd64 go build -o bin/koto-windows-amd64.exe ./cmd/koto
```

### DevContainer環境

このプロジェクトはVS Code / Cursor用のDevContainer環境を含んでいます。

```bash
# VS Code / Cursorで開く
# "Reopen in Container" を選択するだけで開発環境が整います
```

## 🧪 テスト

このプロジェクトは、各層で包括的なテストを実装しています。

```bash
# すべてのテストを実行
go test ./...

# カバレッジレポート
go test -cover ./...

# 詳細なカバレッジレポート
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**テスト統計**:
- Model層: 3テスト関数、7サブテスト
- Repository層: 9テスト関数（インメモリDB使用）
- Service層: 13テスト関数（モックRepository使用）

## 📝 ライセンス

このプロジェクトは[MITライセンス](LICENSE.md)の下で公開されています。

## 🤝 コントリビューション

プルリクエストを歓迎します！バグ報告や機能要望は、[Issues](https://github.com/syeeel/koto-cli-go/issues)にお願いします。

### 開発ガイドライン

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. テストを追加・実行 (`go test ./...`)
4. コミット (`git commit -m 'feat: Add amazing feature'`)
5. プッシュ (`git push origin feature/amazing-feature`)
6. プルリクエストを作成

詳細は[.claude/CLAUDE.md](.claude/CLAUDE.md)の開発ガイドを参照してください。

## 🔗 参考リンク

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUIフレームワーク
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUIコンポーネント集
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - ターミナルスタイリング
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) - Pure Go SQLite


## 💡 FAQ

### Q: データはどこに保存されますか？

A: `~/.koto/koto.db` にSQLiteデータベースとして保存されます。

### Q: 複数のマシンでToDoを同期できますか？

A: 現在、同期機能はありません。`/export` コマンドでJSONファイルにエクスポートし、他のマシンで `/import` することで移行は可能です。

### Q: ToDoの検索機能はありますか？

A: 現在のバージョンでは検索機能は実装されていません。将来のバージョンで追加予定です。

### Q: Windowsで動作しますか？

A: はい、Pure Go実装のため、Windows、macOS、Linuxすべてで動作します。

## 📮 お問い合わせ

バグ報告や質問は、[GitHub Issues](https://github.com/syeeel/koto-cli-go/issues)までお願いします。

---

**koto**を使って、快適なタスク管理を！ 🎵
