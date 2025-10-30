# Changelog

All notable changes to koto will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

## リリース前の変更履歴

v1.0.0以前の開発履歴については、[GitHubのコミット履歴](https://github.com/syeeel/koto-cli-go/commits/main)を参照してください。

### 主要機能（v1.0.0時点）

- ✅ インタラクティブなTUI（Bubbletea）
- ✅ SQLiteによるデータ永続化
- ✅ 優先度管理（高・中・低）
- ✅ 期限管理と期限切れ警告
- ✅ ポモドーロタイマー（25分）
- ✅ 作業時間の自動記録
- ✅ JSON形式でのエクスポート/インポート
- ✅ Vimライクなキーバインド（j/k）
- ✅ ステータスフィルター（未完了/完了済み）
