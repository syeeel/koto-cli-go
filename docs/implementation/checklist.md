# koto - 実装チェックリスト

このチェックリストを使って、実装の進捗を追跡してください。
各項目を完了したら `[ ]` を `[x]` に変更します。

**進捗サマリー**: 0/100 タスク完了

---

## Phase 1: 基盤実装 (0/35)

### 1.1 データモデル (0/7)
- [ ] `internal/model/todo.go` 作成
- [ ] `Todo` 構造体定義
- [ ] `TodoStatus` 型と定数定義（`StatusPending`, `StatusCompleted`）
- [ ] `Priority` 型と定数定義（`PriorityLow`, `PriorityMedium`, `PriorityHigh`）
- [ ] `IsCompleted()` メソッド実装
- [ ] `IsPending()` メソッド実装
- [ ] `IsOverdue()` メソッド実装

#### テスト
- [ ] `internal/model/todo_test.go` 作成
- [ ] `IsCompleted()` のテスト実装
- [ ] `IsOverdue()` のテスト実装

### 1.2 Repository層 (0/18)
- [ ] `internal/repository/repository.go` 作成
- [ ] `TodoRepository` インターフェース定義
- [ ] `migrations/001_init.sql` 作成（スキーマ定義）

#### SQLite実装
- [ ] `internal/repository/sqlite.go` 作成
- [ ] `SQLiteRepository` 構造体定義
- [ ] `NewSQLiteRepository()` 実装
- [ ] `initSchema()` 実装
- [ ] `Create()` 実装
- [ ] `GetByID()` 実装
- [ ] `GetAll()` 実装
- [ ] `GetByStatus()` 実装
- [ ] `Update()` 実装
- [ ] `Delete()` 実装
- [ ] `MarkAsCompleted()` 実装
- [ ] `Close()` 実装

#### テスト
- [ ] `internal/repository/sqlite_test.go` 作成
- [ ] データベース初期化のテスト
- [ ] 各CRUD操作のテスト（最低5つ）

#### 動作確認
- [ ] インメモリDBで動作確認用コードを書いて試す
- [ ] データの作成・取得・更新・削除ができることを確認

### 1.3 Service層 (0/10)
- [ ] `internal/service/todo_service.go` 作成
- [ ] エラー定数定義（`ErrTodoNotFound`, `ErrInvalidTitle`, etc.）
- [ ] `TodoService` 構造体定義
- [ ] `NewTodoService()` 実装
- [ ] `AddTodo()` 実装
- [ ] `EditTodo()` 実装
- [ ] `DeleteTodo()` 実装
- [ ] `CompleteTodo()` 実装
- [ ] `ListTodos()` 実装
- [ ] `ListPendingTodos()` 実装
- [ ] `ListCompletedTodos()` 実装
- [ ] `validateTitle()` 実装

#### テスト
- [ ] `internal/service/todo_service_test.go` 作成
- [ ] バリデーションのテスト
- [ ] 各操作のテスト（モックRepository使用）

---

## Phase 2: MVP実装 (0/30)

### 2.1 TUI基盤 (0/12)
- [ ] `internal/tui/model.go` 作成
- [ ] `ViewMode` 型と定数定義
- [ ] `Model` 構造体定義
- [ ] `NewModel()` 実装
- [ ] `Init()` 実装

#### Update関数
- [ ] `internal/tui/update.go` 作成
- [ ] `Update()` 関数実装
- [ ] Enterキーハンドリング
- [ ] Ctrl+C/Escキーハンドリング
- [ ] Up/Downキーハンドリング
- [ ] `handleEnter()` 実装

#### View関数
- [ ] `internal/tui/views.go` 作成
- [ ] `View()` 関数実装
- [ ] `renderListView()` 実装
- [ ] `renderHelpView()` 実装

#### スタイル
- [ ] `internal/tui/styles.go` 作成
- [ ] lipglossスタイル定義（`titleStyle`, `errorStyle`, etc.）

### 2.2 コマンドパーサー (0/12)
- [ ] `internal/tui/commands.go` 作成
- [ ] `commandExecutedMsg` 型定義
- [ ] `todosLoadedMsg` 型定義
- [ ] `parseAndExecuteCommand()` 実装
- [ ] `loadTodos()` 実装

#### コマンドハンドラー
- [ ] `handleAddCommand()` 実装
- [ ] `handleListCommand()` 実装
- [ ] `handleDoneCommand()` 実装
- [ ] `handleDeleteCommand()` 実装
- [ ] `handleEditCommand()` 実装
- [ ] `handleHelpCommand()` 実装
- [ ] `/quit` 処理実装

### 2.3 メインエントリーポイント (0/6)
- [ ] `internal/config/config.go` 作成
- [ ] データベースパス取得関数実装
- [ ] `.koto` ディレクトリ作成関数実装

#### Main関数
- [ ] `cmd/koto/main.go` 作成
- [ ] データベース初期化処理
- [ ] Serviceインスタンス生成
- [ ] TUIアプリケーション起動
- [ ] エラーハンドリング
- [ ] クリーンアップ処理（defer）

#### 動作確認
- [ ] `go run ./cmd/koto` でアプリが起動する
- [ ] `/add` コマンドでToDoを追加できる
- [ ] `/list` コマンドでToDoが表示される
- [ ] `/done` コマンドでToDoを完了できる
- [ ] `/delete` コマンドでToDoを削除できる
- [ ] アプリを再起動してもデータが残っている

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

- [ ] **M1: 開発環境準備完了** - Phase 0 すべて完了
- [ ] **M2: データ層完成** - Phase 1 すべて完了、テスト通過
- [ ] **M3: MVP完成** - Phase 2 すべて完了、基本操作可能
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
