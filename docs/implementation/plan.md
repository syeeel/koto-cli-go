# koto - 実装計画

## 概要

このドキュメントは、koto（ToDo管理CLIツール）の実装計画を示します。
段階的に機能を実装し、各フェーズで動作確認を行いながら進めます。

## 実装方針

1. **ボトムアップアプローチ**: Model → Repository → Service → TUI の順に実装
2. **段階的リリース**: 各フェーズで動作するバージョンを作成
3. **テスト駆動**: 重要な機能は先にテストを書く
4. **早期フィードバック**: 早い段階でCLIを動かして体験を確認

## フェーズ概要

| フェーズ | 内容 | 目標 | 推定時間 |
|---------|------|------|---------|
| Phase 1 | 基盤実装 | データモデル、Repository、Service | 4-6時間 |
| Phase 2 | MVP実装 | 基本的なTUIと主要コマンド | 6-8時間 |
| Phase 2.5 | 起動バナー表示 | アスキーアートバナーとスプラッシュ画面 | 1-2時間 |
| Phase 3 | 機能拡張 | エクスポート/インポート、UI改善 | 4-6時間 |
| Phase 4 | 品質向上 | テスト、ドキュメント、リリース準備 | 4-6時間 |

**総推定時間**: 20-30時間（2.5-4日相当）


---

## Phase 1: 基盤実装

### 目標
データモデル、Repository、Serviceの実装を完了し、基本的なデータ操作ができる状態にする

### 1.1 データモデル実装 (1-2時間)

#### タスク
- [ ] `internal/model/todo.go` 作成
  - [ ] `Todo` 構造体定義
  - [ ] `TodoStatus` 定数定義
  - [ ] `Priority` 定数定義
  - [ ] ヘルパーメソッド実装（`IsCompleted`, `IsPending`, `IsOverdue`）

**ファイル**: `internal/model/todo.go`

**テスト**: `internal/model/todo_test.go`
- [ ] `IsCompleted()` のテスト
- [ ] `IsOverdue()` のテスト

---

### 1.2 Repository層実装 (2-3時間)

#### タスク
- [ ] `internal/repository/repository.go` 作成
  - [ ] `TodoRepository` インターフェース定義

- [ ] `internal/repository/sqlite.go` 作成
  - [ ] `SQLiteRepository` 構造体定義
  - [ ] `NewSQLiteRepository()` 実装
  - [ ] スキーマ初期化関数 `initSchema()` 実装
  - [ ] CRUD操作の実装
    - [ ] `Create()`
    - [ ] `GetByID()`
    - [ ] `GetAll()`
    - [ ] `GetByStatus()`
    - [ ] `Update()`
    - [ ] `Delete()`
    - [ ] `MarkAsCompleted()`
    - [ ] `Close()`

**ファイル**:
- `internal/repository/repository.go`
- `internal/repository/sqlite.go`
- `migrations/001_init.sql`

**テスト**: `internal/repository/sqlite_test.go`
- [ ] データベース初期化のテスト
- [ ] CRUD操作のテスト（インメモリDB使用）

**検証方法**:
```go
// 簡単な動作確認用のmainを書いて試す
repo, _ := repository.NewSQLiteRepository(":memory:")
todo := &model.Todo{Title: "Test"}
repo.Create(context.Background(), todo)
todos, _ := repo.GetAll(context.Background())
fmt.Println(todos)
```

---

### 1.3 Service層実装 (1-2時間)

#### タスク
- [ ] `internal/service/todo_service.go` 作成
  - [ ] `TodoService` 構造体定義
  - [ ] `NewTodoService()` 実装
  - [ ] エラー定数定義
  - [ ] ビジネスロジック実装
    - [ ] `AddTodo()`
    - [ ] `EditTodo()`
    - [ ] `DeleteTodo()`
    - [ ] `CompleteTodo()`
    - [ ] `ListTodos()`
    - [ ] `ListPendingTodos()`
    - [ ] `ListCompletedTodos()`
    - [ ] `validateTitle()` (private)

**ファイル**: `internal/service/todo_service.go`

**テスト**: `internal/service/todo_service_test.go`
- [ ] バリデーションのテスト
- [ ] 各操作のテスト（モックRepository使用）

**完了条件**: Service層のユニットテストが全てパスすること

---

## Phase 2: MVP実装

### 目標
最小限の機能を持つCLIアプリケーションを完成させる

### 2.1 基本的なTUI実装 (3-4時間)

#### タスク
- [ ] `internal/tui/model.go` 作成
  - [ ] `Model` 構造体定義
  - [ ] `ViewMode` 定数定義
  - [ ] `NewModel()` 実装
  - [ ] `Init()` 実装

- [ ] `internal/tui/update.go` 作成
  - [ ] `Update()` 関数実装
  - [ ] キーボード入力ハンドリング
    - [ ] Enter: コマンド実行
    - [ ] Ctrl+C/Esc: 終了
    - [ ] Up/Down: カーソル移動

- [ ] `internal/tui/views.go` 作成
  - [ ] `View()` 関数実装
  - [ ] `renderListView()` 実装
  - [ ] `renderHelpView()` 実装

- [ ] `internal/tui/styles.go` 作成
  - [ ] lipglossスタイル定義
  - [ ] カラーテーマ設定

**ファイル**:
- `internal/tui/model.go`
- `internal/tui/update.go`
- `internal/tui/views.go`
- `internal/tui/styles.go`

---

### 2.2 コマンドパーサー実装 (2-3時間)

#### タスク
- [ ] `internal/tui/commands.go` 作成
  - [ ] `parseAndExecuteCommand()` 実装
  - [ ] コマンドハンドラー実装
    - [ ] `/add` - `handleAddCommand()`
    - [ ] `/list` - `handleListCommand()`
    - [ ] `/done` - `handleDoneCommand()`
    - [ ] `/delete` - `handleDeleteCommand()`
    - [ ] `/edit` - `handleEditCommand()`
    - [ ] `/help` - `handleHelpCommand()`
    - [ ] `/quit` - 終了処理
  - [ ] メッセージ型定義
    - [ ] `commandExecutedMsg`
    - [ ] `todosLoadedMsg`

**ファイル**: `internal/tui/commands.go`

**コマンド引数パース**:
- [ ] 単純な引数パース（`strings.Fields`）
- [ ] フラグパース（`--desc=`, `--priority=`, `--due=`）

---

### 2.3 メインエントリーポイント実装 (1時間)

#### タスク
- [ ] `cmd/koto/main.go` 作成
  - [ ] データベース初期化
  - [ ] Serviceインスタンス作成
  - [ ] TUIアプリケーション起動
  - [ ] エラーハンドリング
  - [ ] クリーンアップ処理

- [ ] `internal/config/config.go` 作成
  - [ ] データベースパス取得
  - [ ] ディレクトリ作成

**ファイル**:
- `cmd/koto/main.go`
- `internal/config/config.go`

**検証方法**:
```bash
go run ./cmd/koto
# アプリが起動し、/addや/listコマンドが動作することを確認
```

**完了条件**:
- アプリが起動する
- `/add`, `/list`, `/done`, `/delete` が動作する
- データが永続化される（再起動後も残る）

---

## Phase 2.5: 起動バナー表示

### 目標
起動時に「KOTO CLI」のアスキーアートバナーを表示し、ブランドアイデンティティを強化する

### 2.5.1 バナーアスキーアート作成 (30分-1時間)

#### タスク
- [ ] `internal/tui/banner.go` 作成
  - [ ] KOTO CLI アスキーアート定義
  - [ ] バナー文字列を関数として提供
  - [ ] バージョン情報表示関数（オプション）

**ファイル**: `internal/tui/banner.go`

**アスキーアート仕様**:
- 「KOTO CLI」の大きな文字
- ブロック/ピクセル風フォント使用
- 中央揃え
- サブタイトル: "ToDo Manager"

---

### 2.5.2 バナービュー実装 (30分-1時間)

#### タスク
- [ ] `internal/tui/model.go` 更新
  - [ ] `ViewModeBanner` 定数追加
  - [ ] 初期viewModeを `ViewModeBanner` に設定
  - [ ] バナー表示開始時刻フィールド追加（自動遷移用）

- [ ] `internal/tui/views.go` 更新
  - [ ] `renderBannerView()` 関数実装
  - [ ] バナーのスタイリング（lipgloss使用）
  - [ ] "Press any key to continue..." メッセージ

- [ ] `internal/tui/styles.go` 更新
  - [ ] バナー用スタイル定義
  - [ ] カラー選択（ブランドカラー）

- [ ] `internal/tui/update.go` 更新
  - [ ] バナー表示中のキー入力ハンドリング
  - [ ] 任意のキー押下でリスト画面へ遷移
  - [ ] 2秒後の自動遷移（オプション、tea.Tick使用）

**ファイル**:
- `internal/tui/model.go` (更新)
- `internal/tui/views.go` (更新)
- `internal/tui/styles.go` (更新)
- `internal/tui/update.go` (更新)

**検証方法**:
```bash
go build ./cmd/koto
./bin/koto
# バナーが表示され、キー押下でメイン画面に遷移することを確認
```

**完了条件**:
- アプリ起動時にバナーが表示される
- バナーが視覚的に魅力的
- 任意のキーでメイン画面に遷移
- スムーズな画面遷移

---

## Phase 3: 機能拡張

### 目標
エクスポート/インポート機能とUI改善を実装

### 3.1 エクスポート/インポート機能 (2-3時間)

#### タスク
- [ ] Service層にエクスポート/インポート追加
  - [ ] `ExportToJSON()` 実装
  - [ ] `ImportFromJSON()` 実装

- [ ] TUI層にコマンド追加
  - [ ] `/export` コマンドハンドラー実装
  - [ ] `/import` コマンドハンドラー実装
  - [ ] ファイルパス指定の処理
  - [ ] インポート時の確認ダイアログ

**ファイル**:
- `internal/service/todo_service.go` (既存に追加)
- `internal/tui/commands.go` (既存に追加)

**テスト**:
- [ ] JSONエクスポートのテスト
- [ ] JSONインポートのテスト
- [ ] 不正なJSONの処理テスト

---

### 3.2 UI/UX改善 (2-3時間)

#### タスク
- [ ] リスト表示の改善
  - [ ] 優先度の視覚化（絵文字/色）
  - [ ] 期限の表示
  - [ ] 期限切れの強調表示
  - [ ] ステータスの視覚化

- [ ] 詳細表示モード追加
  - [ ] ToDoの詳細を表示する画面
  - [ ] 説明、作成日時、更新日時の表示

- [ ] 編集モードの改善
  - [ ] インタラクティブな入力フォーム
  - [ ] 現在の値をデフォルト表示

- [ ] ヘルプの充実
  - [ ] コマンド一覧
  - [ ] 使用例
  - [ ] キーボードショートカット

**ファイル**:
- `internal/tui/views.go` (拡張)
- `internal/tui/styles.go` (拡張)

---

## Phase 4: 品質向上とリリース準備

### 目標
テスト、ドキュメント、ビルド設定を整え、リリース可能な状態にする

### 4.1 テスト拡充 (2-3時間)

#### タスク
- [ ] ユニットテストの拡充
  - [ ] Model層のテスト
  - [ ] Repository層のテスト
  - [ ] Service層のテスト

- [ ] 統合テストの作成
  - [ ] Repository + Service の統合テスト
  - [ ] エンドツーエンドのシナリオテスト

- [ ] テストカバレッジ確認
  ```bash
  go test -cover ./...
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

**目標カバレッジ**: 80%以上

---

### 4.2 ドキュメント整備 (1-2時間)

#### タスク
- [ ] `README.md` 作成
  - [ ] プロジェクト概要
  - [ ] インストール方法
  - [ ] 使い方
  - [ ] コマンドリファレンス
  - [ ] スクリーンショット/デモGIF

- [ ] `CONTRIBUTING.md` 作成（オプション）
  - [ ] 開発環境のセットアップ
  - [ ] コーディング規約
  - [ ] プルリクエストの手順

- [ ] コード内コメントの追加
  - [ ] 公開API用のGoDoc形式コメント
  - [ ] 複雑なロジックの説明

---

### 4.3 ビルドとリリース設定 (1-2時間)

#### タスク
- [ ] Makefile最終調整
  - [ ] `make build`
  - [ ] `make test`
  - [ ] `make lint`
  - [ ] `make clean`
  - [ ] `make install`

- [ ] GoReleaser設定
  - [ ] `.goreleaser.yml` 作成
  - [ ] ビルドターゲット設定
  - [ ] アーカイブ設定
  - [ ] チェックサム生成

- [ ] GitHub Actions設定（オプション）
  - [ ] `.github/workflows/test.yml` (CI)
  - [ ] `.github/workflows/release.yml` (リリース自動化)

- [ ] バージョン管理
  - [ ] セマンティックバージョニング採用
  - [ ] `CHANGELOG.md` 作成

**ファイル**:
- `Makefile`
- `.goreleaser.yml`
- `.github/workflows/test.yml`
- `.github/workflows/release.yml`
- `CHANGELOG.md`

---

### 4.4 初回リリース (30分-1時間)

#### タスク
- [ ] バージョン v1.0.0 のタグ作成
- [ ] GitHub Releasesで公開
  - [ ] リリースノート作成
  - [ ] バイナリのアップロード
  - [ ] チェックサム確認

- [ ] インストールテスト
  - [ ] 各プラットフォームでバイナリが動作することを確認
  - [ ] `go install` での インストールテスト

**検証**:
```bash
# Linux
wget https://github.com/user/koto/releases/download/v1.0.0/koto-linux-amd64
chmod +x koto-linux-amd64
./koto-linux-amd64

# go install
go install github.com/user/koto@latest
koto
```

---

## 優先度付きタスクリスト

### 必須（MVP）
1. Phase 0: プロジェクトセットアップ
2. Phase 1: 基盤実装
3. Phase 2: MVP実装
4. Phase 4.3: 基本的なビルド設定

### 重要（v1.0）
5. Phase 3.1: エクスポート/インポート機能
6. Phase 4.1: テスト拡充
7. Phase 4.2: ドキュメント整備
8. Phase 4.4: 初回リリース

### あると良い（v1.1+）
9. Phase 3.2: UI/UX改善
10. GitHub Actions設定
11. CONTRIBUTING.md

---

## リスクと対策

### リスク1: modernc.org/sqlite の動作問題
**影響**: データ永続化ができない
**対策**:
- 早い段階でSQLite接続テストを実施
- 問題があれば go-sqlite3 への切り替えも検討

### リスク2: Bubbletea の学習曲線
**影響**: UI実装に時間がかかる
**対策**:
- 公式サンプルを参考にする
- シンプルなUIから始める
- examples/ フォルダを参照

### リスク3: スコープクリープ
**影響**: 開発が長引く
**対策**:
- MVPに集中する
- 追加機能は v2.0 以降に回す
- Phase単位で完了判定を厳格に行う

---

## マイルストーン

| マイルストーン | 完了条件 | 期日 |
|--------------|---------|------|
| M1: 開発環境準備完了 | Phase 0 完了 | Day 1 |
| M2: データ層完成 | Phase 1 完了、テスト通過 | Day 1-2 |
| M3: MVP完成 | Phase 2 完了、基本操作可能 | Day 2-3 |
| M4: v1.0リリース | Phase 3, 4 完了 | Day 3-4 |

---

## 参考リソース

### Bubbletea
- [公式ドキュメント](https://github.com/charmbracelet/bubbletea)
- [Bubblesコンポーネント](https://github.com/charmbracelet/bubbles)
- [サンプルアプリ](https://github.com/charmbracelet/bubbletea/tree/master/examples)

### modernc.org/sqlite
- [ドキュメント](https://pkg.go.dev/modernc.org/sqlite)
- [使用例](https://gitlab.com/cznic/sqlite/-/blob/master/README.md)

### Go Testing
- [Go Testing Guide](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)

---

## 次のステップ

1. このplan.mdとchecklist.mdを確認
2. Phase 0から開始
3. 各フェーズ完了後、checklist.mdを更新
4. 問題が発生したら、このドキュメントを更新

**実装を開始する準備ができました！**
