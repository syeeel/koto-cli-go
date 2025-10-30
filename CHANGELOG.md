# Changelog

All notable changes to koto will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.6] - 2025-10-30

### Added
- **レスポンシブレイアウト**: ターミナル幅に応じて動的にレイアウトを調整
- **最小ターミナル幅チェック**: 100文字未満の場合にエラーメッセージを表示
- **動的幅計算システム**: 全画面の幅を一元管理する`DynamicWidths`構造体を追加
- **クロスターミナル互換ASCIIボーダー**: シンプルなASCII文字（+, -, |）を使用したカスタムボーダー
- **動的バージョン情報システム**: ビルド時に自動的にバージョン、コミットハッシュ、ビルド日時を注入
- **Makefile**: 開発時のビルドを効率化する包括的なMakefileを追加（build, test, clean, install, run, release等）

### Fixed
- **タスク一覧**: フォーカス時に選択行が2行表示される問題を修正
- **起動画面**: 枠線が途切れる問題を修正（ASCIIアートとToDoボックスを中央配置）
- **詳細画面**: 枠線が途切れる問題を修正（全ボックスを動的幅に対応）
- **ポモドーロ画面**: 枠線が途切れる問題を修正（プログレスバー、情報ボックスを動的幅に対応）
- **ポモドーロ画面**: タスク情報（Task IDとタスク名）を中央表示に変更
- **ポモドーロ画面**: プログレスバーを中央表示に改善
- **macOS Terminal互換性**: Unicode box-drawing charactersをASCII文字に置き換え、全てのターミナルで正しく表示されるように修正

### Changed
- **メインリスト画面**: カラム幅を動的に計算、タイトルカラムを可変幅に変更
- **詳細画面**: 3カラムレイアウトを比例配分で動的調整
- **ポモドーロ画面**: 全要素をターミナル幅に応じてセンタリング
- **ボーダースタイル**: RoundedBorder/NormalBorderから互換性の高いASCIIボーダーに変更
- **バージョン表示**: 起動画面とコマンドラインで詳細なバージョン情報（コミット、ビルド日時）を表示

### Technical
- `internal/tui/styles.go`: 動的幅計算関数とヘルパー関数を追加、ASCIIボーダー定義を追加
- `internal/tui/views.go`: 全画面のレスポンシブ化、最小幅チェック機能を追加、全ボーダーをASCIIに置換、プログレスバー中央表示改善
- `internal/tui/banner.go`: バージョン情報変数をエクスポート、`GetVersion()`を詳細表示に変更
- `cmd/koto/main.go`: TUIパッケージにバージョン情報を注入
- `Makefile`: ビルド時にldflagsでバージョン情報を設定、各種開発タスクを自動化
- Lipglossの`PlaceHorizontal`を活用した中央配置の実装
- 全てのUnicode box-drawing charactersをASCII文字に置換してクロスターミナル互換性を向上
- GoReleaserとの統合により、リリースビルドで自動的にバージョン情報が注入される

## [1.0.3] - 2025-10-30

### Fixed
- **install.sh**: アーカイブ名の形式を修正（koto-cli-goプロジェクト名を使用）
- **install.sh**: バージョン番号から'v'プレフィックスを削除してGoReleaserの実際のアーカイブ名と一致

### Changed
- インストールスクリプトがmacOS/Linuxで正常に動作するように修正

## [1.0.2] - 2025-10-30

### Changed
- Homebrew Tap設定を一時的に無効化（セットアップ準備中）
- GoReleaser v2互換性のため非推奨設定（format_overrides）を削除
- アーカイブ設定を簡素化

### Fixed
- リリース時のHomebrew Tap 401エラーを修正（設定を無効化）

## [1.0.1] - 2025-10-30

### Changed
- READMEのHomebrew案内を「準備中🚧」ステータスに更新
- 次回リリース（v1.0.1以降）からHomebrewが利用可能になることを明記

## [1.0.0] - 2025-10-30

### Added
- **GoReleaser統合**: 自動リリースワークフロー
- **マルチプラットフォームビルド**: macOS (Intel/Apple Silicon), Linux (amd64/arm64), Windows (amd64)
- **インストールスクリプト**: ワンラインインストール（curl | sh）
- **バージョン表示**: `--version`フラグでバージョン、コミット、ビルド日時を表示
- **GitHub Actions**: タグプッシュ時の自動リリースワークフロー

### Documentation
- リリースガイド（docs/RELEASE.md）
- Homebrewセットアップガイド（docs/SETUP_HOMEBREW.md）
- クイックスタートガイド（docs/QUICKSTART_RELEASE.md）

### Infrastructure
- `.goreleaser.yaml`: GoReleaser設定
- `.github/workflows/release.yml`: 自動リリースワークフロー
- `install.sh`: ユーザー向けインストールスクリプト

---