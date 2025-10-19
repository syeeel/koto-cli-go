# koto - 実装チェックリスト

このチェックリストを使って、実装の進捗を追跡してください。
各項目を完了したら `[ ]` を `[x]` に変更します。

**進捗サマリー**: 71/106 タスク完了

---

## Phase 1: 基盤実装 (35/35) ✅ COMPLETE

### 1.1 データモデル (10/10) ✅
- [x] `internal/model/todo.go` 作成
- [x] `Todo` 構造体定義
- [x] `TodoStatus` 型と定数定義（`StatusPending`, `StatusCompleted`）
- [x] `Priority` 型と定数定義（`PriorityLow`, `PriorityMedium`, `PriorityHigh`）
- [x] `IsCompleted()` メソッド実装
- [x] `IsPending()` メソッド実装
- [x] `IsOverdue()` メソッド実装

#### テスト
- [x] `internal/model/todo_test.go` 作成
- [x] `IsCompleted()` のテスト実装
- [x] `IsOverdue()` のテスト実装

### 1.2 Repository層 (18/18) ✅
- [x] `internal/repository/repository.go` 作成
- [x] `TodoRepository` インターフェース定義
- [x] `migrations/001_init.sql` 作成（スキーマ定義）

#### SQLite実装
- [x] `internal/repository/sqlite.go` 作成
- [x] `SQLiteRepository` 構造体定義
- [x] `NewSQLiteRepository()` 実装
- [x] `initSchema()` 実装
- [x] `Create()` 実装
- [x] `GetByID()` 実装
- [x] `GetAll()` 実装
- [x] `GetByStatus()` 実装
- [x] `Update()` 実装
- [x] `Delete()` 実装
- [x] `MarkAsCompleted()` 実装
- [x] `Close()` 実装

#### テスト
- [x] `internal/repository/sqlite_test.go` 作成
- [x] データベース初期化のテスト
- [x] 各CRUD操作のテスト（最低5つ）

#### 動作確認
- [x] インメモリDBで動作確認用コードを書いて試す
- [x] データの作成・取得・更新・削除ができることを確認

### 1.3 Service層 (15/15) ✅
- [x] `internal/service/todo_service.go` 作成
- [x] エラー定数定義（`ErrTodoNotFound`, `ErrInvalidTitle`, etc.）
- [x] `TodoService` 構造体定義
- [x] `NewTodoService()` 実装
- [x] `AddTodo()` 実装
- [x] `EditTodo()` 実装
- [x] `DeleteTodo()` 実装
- [x] `CompleteTodo()` 実装
- [x] `ListTodos()` 実装
- [x] `ListPendingTodos()` 実装
- [x] `ListCompletedTodos()` 実装
- [x] `validateTitle()` 実装
- [x] `ExportToJSON()` 実装
- [x] `ImportFromJSON()` 実装
- [x] `validatePriority()` 実装

#### テスト
- [x] `internal/service/todo_service_test.go` 作成
- [x] バリデーションのテスト
- [x] 各操作のテスト（モックRepository使用）

---

## Phase 2: MVP実装 (30/30) ✅ COMPLETE

### 2.1 TUI基盤 (12/12) ✅
- [x] `internal/tui/model.go` 作成
- [x] `ViewMode` 型と定数定義
- [x] `Model` 構造体定義
- [x] `NewModel()` 実装
- [x] `Init()` 実装

#### Update関数
- [x] `internal/tui/update.go` 作成
- [x] `Update()` 関数実装
- [x] Enterキーハンドリング
- [x] Ctrl+C/Escキーハンドリング
- [x] Up/Downキーハンドリング
- [x] `handleEnter()` 実装

#### View関数
- [x] `internal/tui/views.go` 作成
- [x] `View()` 関数実装
- [x] `renderListView()` 実装
- [x] `renderHelpView()` 実装

#### スタイル
- [x] `internal/tui/styles.go` 作成
- [x] lipglossスタイル定義（`titleStyle`, `errorStyle`, etc.）

### 2.2 コマンドパーサー (12/12) ✅
- [x] `internal/tui/commands.go` 作成
- [x] `commandExecutedMsg` 型定義
- [x] `todosLoadedMsg` 型定義
- [x] `parseAndExecuteCommand()` 実装
- [x] `loadTodos()` 実装

#### コマンドハンドラー
- [x] `handleAddCommand()` 実装
- [x] `handleListCommand()` 実装
- [x] `handleDoneCommand()` 実装
- [x] `handleDeleteCommand()` 実装
- [x] `handleEditCommand()` 実装
- [x] `handleHelpCommand()` 実装
- [x] `/quit` 処理実装

### 2.3 メインエントリーポイント (6/6) ✅
- [x] `internal/config/config.go` 作成
- [x] データベースパス取得関数実装
- [x] `.koto` ディレクトリ作成関数実装

#### Main関数
- [x] `cmd/koto/main.go` 作成
- [x] データベース初期化処理
- [x] Serviceインスタンス生成
- [x] TUIアプリケーション起動
- [x] エラーハンドリング
- [x] クリーンアップ処理（defer）

#### 動作確認
- [x] ネットワーク接続問題を診断・解決（devcontainer firewall無効化）
- [x] go mod download成功
- [x] go mod tidy成功（go.sum生成）
- [x] 全テスト通過（Model, Repository, Service層）
- [x] `go build ./cmd/koto` 成功（bin/koto 9.8MB）
- [x] データベース初期化確認（~/.koto/koto.db作成）
- [x] スキーマ検証（テーブル・インデックス確認）
- [ ] インタラクティブターミナルでの完全動作確認（要ユーザー実行）
  - [ ] `/add` コマンドでToDoを追加できる
  - [ ] `/list` コマンドでToDoが表示される
  - [ ] `/done` コマンドでToDoを完了できる
  - [ ] `/delete` コマンドでToDoを削除できる
  - [ ] アプリを再起動してもデータが残っている

---

## Phase 2.5: 起動バナー表示 (6/6) ✅ COMPLETE

### 2.5.1 バナーアスキーアート (2/2) ✅
- [x] `internal/tui/banner.go` 作成
- [x] KOTO CLI アスキーアート定義

### 2.5.2 バナービュー実装 (4/4) ✅
- [x] `ViewModeBanner` を model.go に追加
- [x] `renderBannerView()` を views.go に実装
- [x] バナー用スタイルを styles.go に追加
- [x] バナー表示からリスト画面への遷移を update.go に実装
  - [x] 任意のキー押下で遷移

#### 動作確認
- [x] ビルド成功（bin/koto 9.8MB）
- [x] 全テスト通過
- [ ] 起動時にバナーが表示される（要ユーザー実行）
- [ ] キー押下でメイン画面に遷移する（要ユーザー実行）

---

## Phase 3: 機能拡張 (0/12)

### 3.1 エクスポート/インポート (0/8)
#### Service層
- [ ] `ExportToJSON()` メソッド追加
- [ ] `ImportFromJSON()` メソッド追加
- [ ] エクスポート用エラー定数追加
- [ ] インポート用エラー定数追加

#### TUI層
- [ ] `/export` コマンドハンドラー実装
- [ ] `/import` コマンドハンドラー実装
- [ ] ファイルパス指定の処理
- [ ] インポート時の確認ダイアログ実装

#### テスト
- [ ] エクスポート機能のテスト
- [ ] インポート機能のテスト
- [ ] 不正なJSONファイルの処理テスト

#### 動作確認
- [ ] `/export` でJSONファイルが作成される
- [ ] `/import` でJSONファイルを読み込める

### 3.2 UI/UX改善 (0/4)
- [ ] 優先度の視覚化（絵文字/色）
- [ ] 期限の表示
- [ ] 期限切れの強調表示
- [ ] ステータスの視覚化改善

#### 詳細表示
- [ ] ToDo詳細表示モード追加
- [ ] 説明、作成日時、更新日時の表示

#### ヘルプ充実
- [ ] コマンド一覧を `/help` に追加
- [ ] 使用例を追加
- [ ] キーボードショートカット一覧

---

## Phase 4: 品質向上とリリース (0/13)

### 4.1 テスト拡充 (0/4)
- [ ] Model層のテストカバレッジ確認
- [ ] Repository層のテストカバレッジ確認
- [ ] Service層のテストカバレッジ確認
- [ ] 統合テスト作成

#### カバレッジ確認
- [ ] `go test -cover ./...` 実行
- [ ] カバレッジ80%以上を確認
- [ ] 不足している部分のテスト追加

### 4.2 ドキュメント整備 (0/4)
- [ ] `README.md` 作成
  - [ ] プロジェクト概要
  - [ ] インストール方法
  - [ ] 使い方
  - [ ] コマンドリファレンス
  - [ ] スクリーンショット/デモGIF（オプション）

- [ ] コード内コメント追加
  - [ ] 公開APIのGoDocコメント
  - [ ] 複雑なロジックの説明コメント

### 4.3 ビルドとリリース設定 (0/5)
- [ ] `Makefile` 最終調整
  - [ ] `make build` 動作確認
  - [ ] `make test` 動作確認
  - [ ] `make clean` 動作確認
  - [ ] `make build-all` でクロスコンパイル確認

- [ ] `.goreleaser.yml` 作成
- [ ] `CHANGELOG.md` 作成

#### GitHub Actions（オプション）
- [ ] `.github/workflows/test.yml` 作成
- [ ] `.github/workflows/release.yml` 作成

### 4.4 初回リリース (0/4)
- [ ] バージョン v1.0.0 のタグ作成
- [ ] GitHub Releases作成
- [ ] リリースノート作成
- [ ] バイナリのアップロード

#### 検証
- [ ] Linux版バイナリの動作確認
- [ ] macOS版バイナリの動作確認（可能であれば）
- [ ] Windows版バイナリの動作確認（可能であれば）
- [ ] `go install` でのインストールテスト

---

## オプショナルタスク

### 追加機能（v1.1以降）
- [ ] フィルタリング機能（優先度、期限）
- [ ] ソート機能
- [ ] 検索機能
- [ ] 統計表示（完了率など）

### インフラ・ツール
- [ ] GitHub Discussions設定
- [ ] Issue テンプレート作成
- [ ] PR テンプレート作成
- [ ] `CONTRIBUTING.md` 作成

### ドキュメント拡張
- [ ] ユーザーガイド作成
- [ ] アーキテクチャ図作成
- [ ] FAQ作成

---

## マイルストーン達成チェック

- [x] **M1: 開発環境準備完了** - Phase 0 すべて完了 ✅
- [x] **M2: データ層完成** - Phase 1 すべて完了、テスト通過 ✅
- [x] **M3: MVP実装完了** - Phase 2 コード実装完了 ✅ (動作確認はネットワーク接続後)
- [ ] **M4: v1.0リリース** - Phase 3, 4 すべて完了

---

## 完了時の最終確認

リリース前に、以下をすべて確認してください：

### 機能
- [ ] すべての基本コマンドが動作する（/add, /edit, /delete, /done, /list）
- [ ] エクスポート/インポートが動作する
- [ ] データが正しく永続化される
- [ ] エラーが適切に処理される

### コード品質
- [ ] テストカバレッジ80%以上
- [ ] `go vet ./...` がエラーなし
- [ ] `golangci-lint run` がエラーなし（設定している場合）
- [ ] すべてのテストがパス

### ドキュメント
- [ ] README.mdが完成している
- [ ] 設計ドキュメントが最新
- [ ] CHANGELOGが記載されている

### ビルド
- [ ] `make build` が成功する
- [ ] クロスコンパイルが成功する
- [ ] バイナリサイズが妥当（目安: 10-20MB）

### リリース
- [ ] GitHub Releasesページが作成されている
- [ ] バイナリがアップロードされている
- [ ] リリースノートが記載されている

---

**実装頑張ってください！各タスクを完了するたびに達成感を味わいましょう！**
