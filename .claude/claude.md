# koto プロジェクト - Claude Code 開発ガイド

## プロジェクト概要

**koto**は、Goで開発されたインタラクティブなToDoリスト管理CLIツールです。
bubbletea フレームワークを使用して、リッチで直感的なターミナルUIを提供します。

### 技術スタック
- **言語**: Go 1.21+
- **TUIフレームワーク**: Bubbletea, Bubbles, Lipgloss
- **データベース**: SQLite (modernc.org/sqlite)
- **アーキテクチャ**: レイヤードアーキテクチャ (Model → Repository → Service → TUI)

---

## MCP (Model Context Protocol) サーバー

このプロジェクトでは以下のMCPサーバーが設定されています（`.claude/mcp.json`）：

### filesystem
- ワークスペース (`/workspace`) とkoto設定ディレクトリ (`~/.koto`) へのアクセス
- ファイル操作の効率化に使用

### sqlite
- テストデータベース (`/workspace/test.db`) へのアクセス
- データベース操作やクエリ実行に使用

### context7
- Upstashのセマンティック検索/RAG機能
- コードベースの理解や検索の効率化に使用

**注意**: MCPサーバーはClaude Codeが自動的に活用します。開発者が直接操作する必要はありません。

---

## 重要なドキュメント

開発を開始する前に、必ず以下のドキュメントを参照してください：

### 設計書
- **@docs/design/basic_design.md** - 基本設計書（プロジェクト概要、主要機能、アーキテクチャ）
- **@docs/design/detailed_design.md** - 詳細設計書（各層の詳細仕様、データモデル、API設計）

### 実装管理
- **@docs/implementation/plan.md** - 実装計画（フェーズ別タスク、推定時間、マイルストーン）
- **@docs/implementation/checklist.md** - 実装チェックリスト（進捗管理、タスク一覧）
- **@docs/implementation/log.md** - 実装ログ（実装履歴、変更内容、課題、決定事項）

---

## 開発ワークフロー

### 1. タスク開始前

1. **@docs/implementation/plan.md** で現在のフェーズを確認
2. **@docs/implementation/checklist.md** で次に実装すべきタスクを確認
3. **@docs/design/detailed_design.md** で該当する詳細仕様を確認
4. **@docs/implementation/log.md** で過去の実装履歴や決定事項を確認

### 2. 実装中

1. **TodoWrite** ツールを使用してタスクを追跡
2. 設計書に従って実装
3. テストを書く（TDD推奨）
4. コードレビュー用に **golang-pro** サブエージェントを使用

### 3. タスク完了後

1. **@docs/implementation/checklist.md** の該当項目を `[x]` に更新
2. **@docs/implementation/log.md** に実装内容を記録：
   ```markdown
   ## YYYY-MM-DD HH:MM - [完了したタスク名]

   ### 実装内容
   - 実装した機能の詳細
   - 変更したファイル

---

## サブエージェントの活用

### golang-pro エージェント

**コードを書く前とコードを書いた後に必ず使用してください。**

#### 使用タイミング
- **実装前**: テスト計画とアーキテクチャレビュー
- **実装後**: コードレビュー、ベストプラクティスチェック

#### 使用例
```
実装前:
- "次のタスクについて、テスト計画とアーキテクチャを golang-pro でレビュー"

実装後:
- "実装したコードを golang-pro でレビューして改善提案をもらう"
```

### Explore エージェント

**コードベースの探索や理解に使用してください。**

#### 使用タイミング
- 既存のコード構造を理解したい時
- 特定の機能がどこに実装されているか探す時
- アーキテクチャパターンを調査する時

#### thoroughnessレベル
- `quick`: 基本的な検索（ファイル名やキーワード）
- `medium`: 中程度の探索（複数パターン）
- `very thorough`: 徹底的な分析（アーキテクチャ全体）

---

## コーディング規約

### Goのベストプラクティス

1. **エラーハンドリング**
   - すべてのエラーを適切に処理
   - カスタムエラーは `errors.New()` または `fmt.Errorf()` で作成
   - センチネルエラーは定数として定義

2. **命名規則**
   - パッケージ名: 小文字、単一単語
   - エクスポート: PascalCase
   - 非エクスポート: camelCase
   - インターフェース: `-er` サフィックス（例: `TodoRepository`）

3. **テスト**
   - ファイル名: `*_test.go`
   - テーブル駆動テスト推奨
   - モックは必要に応じて作成
   - カバレッジ目標: 80%以上

4. **コメント**
   - エクスポートされる要素にはGoDoc形式のコメント必須
   - 複雑なロジックには説明コメント
   - TODO/FIXME/NOTE を適切に使用

### プロジェクト固有のルール

1. **レイヤーの分離**
   - Model層: ビジネスロジックなし、純粋なデータ構造
   - Repository層: データアクセスのみ、ビジネスロジックなし
   - Service層: ビジネスロジック、バリデーション
   - TUI層: UI/UX、ユーザー入力処理

2. **依存関係**
   - 上位層から下位層への依存のみ許可
   - TUI → Service → Repository → Model
   - 逆方向の依存は禁止

3. **コンテキスト**
   - Repository層のすべてのメソッドは `context.Context` を第一引数に取る
   - タイムアウトやキャンセル処理を適切に実装

---

## 実装ログの記録方法

すべての重要な実装作業の後、**@docs/implementation/log.md** に記録してください。

### ログ記録のテンプレート

```markdown
## YYYY-MM-DD HH:MM - [タスク名]

### 実装内容
- 箇条書きで実装した内容
- 追加/変更したファイルのリスト

### 技術的な決定
- なぜその実装方法を選んだか
- 検討した代替案とその却下理由
- パフォーマンスやセキュリティの考慮事項

### 遭遇した問題と解決方法
- 問題の詳細
- 解決策
- 学んだこと

### テスト
- 実装したテストの概要
- カバレッジ状況

### 次のステップ
- 残っている課題
- 次に取り組むタスク
- 改善案
```

### ログ記録すべき内容

1. **必須**
   - 新機能の実装
   - 重要なバグ修正
   - アーキテクチャの変更
   - 設計の決定や変更

2. **推奨**
   - 困難な問題の解決
   - パフォーマンスの改善
   - 依存関係の追加や更新
   - テストカバレッジの向上

3. **オプション**
   - 小さなリファクタリング
   - コメントの追加
   - ドキュメントの更新

---

## チェックリストの更新

**@docs/implementation/checklist.md** は進捗管理の重要なツールです。

### 更新ルール

1. **タスク完了時**
   - `[ ]` を `[x]` に変更
   - 進捗サマリーの数字を更新（例: `0/100` → `1/100`）

2. **新しいタスクを発見した時**
   - 適切なフェーズに追加
   - 進捗サマリーの総数を更新

3. **タスクが不要になった時**
   - 削除せず、打ち消し線（`~~タスク名~~`）を使用
   - 理由を log.md に記録

### 更新のタイミング

- タスク完了後、すぐに更新
- 1日の作業終了時に見直し
- フェーズ完了時に次フェーズを確認

---

## デバッグとトラブルシューティング

### よくある問題

1. **SQLite接続エラー**
   - データベースファイルのパーミッション確認
   - パスが正しいか確認
   - インメモリDB（`:memory:`）でテスト

2. **Bubbletea UIの問題**
   - `Update()` 関数の戻り値を確認
   - メッセージの型が正しいか確認
   - 公式サンプルを参照

3. **テストの失敗**
   - モックが正しく設定されているか
   - テストデータの準備を確認
   - `go test -v` で詳細を確認

### デバッグ手順

1. **ログを確認**
   - `fmt.Printf()` でデバッグ出力（TUIでは使用注意）
   - エラーメッセージを読む

2. **テストを実行**
   ```bash
   go test -v ./...
   go test -run TestSpecificFunction
   ```

3. **サブエージェントに相談**
   - golang-pro エージェントでコードレビュー
   - Explore エージェントで類似コードを探す

---

## Git コミットガイドライン

### CHANGELOG.mdの更新
git addを行う前に、**@CHANGELOG.md**に変更履歴を記載してください。

### コミットメッセージ

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type
- `feat`: 新機能
- `fix`: バグ修正
- `refactor`: リファクタリング
- `test`: テスト追加・修正
- `docs`: ドキュメント
- `chore`: ビルド、ツール設定

### 例
```
feat(service): Add TodoService with basic CRUD operations

- Implement AddTodo, EditTodo, DeleteTodo
- Add validation for todo title
- Include unit tests with 85% coverage

Refs: #12
```

---

## リリース前チェックリスト

Phase 4完了時に、以下を確認してください：

### コード品質
- [ ] `go vet ./...` がエラーなし
- [ ] `go test ./...` が全てパス
- [ ] テストカバレッジ 80%以上
- [ ] GoDocコメントが完備

### 機能
- [ ] すべての基本コマンドが動作
- [ ] データが正しく永続化される
- [ ] エラーハンドリングが適切

### ドキュメント
- [ ] README.md が完成
- [ ] 設計書が最新
- [ ] CHANGELOG.md が記載

### ビルド
- [ ] `make build` が成功
- [ ] クロスコンパイルが成功
- [ ] バイナリサイズが妥当（10-20MB）

---

## 参考リソース

### Bubbletea
- [公式ドキュメント](https://github.com/charmbracelet/bubbletea)
- [Bubbles コンポーネント](https://github.com/charmbracelet/bubbles)
- [サンプルアプリ](https://github.com/charmbracelet/bubbletea/tree/master/examples)

### Go
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)

### SQLite
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)

---

## 質問やヘルプが必要な時

1. **設計に関する質問**
   - @docs/design/basic_design.md を確認
   - @docs/design/detailed_design.md で詳細仕様を確認

2. **実装に関する質問**
   - @docs/implementation/plan.md でフェーズを確認
   - golang-pro サブエージェントに相談

3. **過去の決定を確認したい**
   - @docs/implementation/log.md を検索

4. **進捗を確認したい**
   - @docs/implementation/checklist.md を参照

---
