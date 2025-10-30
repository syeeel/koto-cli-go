# クイックスタート: 初回リリース

このガイドでは、kotoの初回リリースを最短手順で実行する方法を説明します。

## 📋 事前チェック

リリース前に以下を確認してください：

```bash
# テストが通るか確認
go test ./...

# ビルドが成功するか確認
go build -o bin/koto ./cmd/koto

# GoReleaserの設定が正しいか確認
goreleaser check

# ローカルビルドテスト
goreleaser build --snapshot --clean --single-target
./dist/koto_*/koto --version
```

## 🍺 Homebrew Tap のセットアップ（初回のみ）

### ステップ 1: リポジトリ作成

```bash
gh repo create syeeel/homebrew-tap --public --description "Homebrew tap for koto CLI"
```

または [GitHub UI](https://github.com/new) で作成

### ステップ 2: Personal Access Token 作成

1. [Personal Access Tokens](https://github.com/settings/tokens) ページを開く
2. "Generate new token (classic)" をクリック
3. スコープ: `public_repo` ✅
4. トークンをコピー

### ステップ 3: GitHub Secrets に追加

1. [Settings → Secrets](https://github.com/syeeel/koto-cli-go/settings/secrets/actions) を開く
2. "New repository secret" をクリック
3. Name: `TAP_GITHUB_TOKEN`
4. Secret: コピーしたトークン
5. "Add secret" をクリック

詳細手順: [docs/SETUP_HOMEBREW.md](SETUP_HOMEBREW.md)

## 🚀 初回リリース実行

### ステップ 1: 変更をコミット

```bash
# 現在の変更を確認
git status

# すべての変更をステージング
git add .

# コミット
git commit -m "feat: Add GoReleaser setup and Homebrew support"

# メインブランチにプッシュ
git push origin main
```

### ステップ 2: リリースタグを作成

```bash
# タグを作成（バージョンは適宜変更）
git tag -a v0.1.0 -m "Initial release

- Interactive TUI for todo management
- SQLite database for persistence
- Priority and due date support
- Pomodoro timer integration
- Homebrew tap support
"

# タグをプッシュ（これでリリースが自動開始）
git push origin v0.1.0
```

### ステップ 3: GitHub Actions を確認

1. [Actions タブ](https://github.com/syeeel/koto-cli-go/actions) を開く
2. "Release" ワークフローの進行を確認
3. 完了まで約2-3分待つ

ワークフローの流れ：
```
テスト実行
  ↓
ビルド（各プラットフォーム）
  ↓
GitHub Releases に公開
  ↓
Homebrew Formula 更新
  ↓
完了 ✅
```

### ステップ 4: リリースの確認

#### GitHub Releases

[Releases ページ](https://github.com/syeeel/koto-cli-go/releases) で以下を確認：

- ✅ 新しいリリース `v0.1.0` が作成されている
- ✅ 各プラットフォーム向けバイナリがアップロードされている
  - `koto_0.1.0_darwin_amd64.tar.gz`
  - `koto_0.1.0_darwin_arm64.tar.gz`
  - `koto_0.1.0_linux_amd64.tar.gz`
  - `koto_0.1.0_linux_arm64.tar.gz`
  - `koto_0.1.0_windows_amd64.zip`
  - `checksums.txt`
- ✅ Changelog が自動生成されている

#### Homebrew Tap

[homebrew-tap リポジトリ](https://github.com/syeeel/homebrew-tap) で以下を確認：

- ✅ `Formula/koto.rb` ファイルが作成されている
- ✅ ファイル内にバージョン `0.1.0` が記載されている
- ✅ ダウンロードURL と SHA256 ハッシュが記載されている

## 🧪 インストールのテスト

### インストールスクリプト（macOS/Linux）

```bash
curl -sSfL https://raw.githubusercontent.com/syeeel/koto-cli-go/main/install.sh | sh
```

### Homebrew（macOS）

```bash
brew tap syeeel/tap
brew install koto
```

### 動作確認

```bash
# バージョン確認
koto --version
# 出力例: koto version 0.1.0

# アプリ起動
koto
```

## 📣 リリースのアナウンス

リリースが成功したら、以下でアナウンスしましょう：

### README.md を更新

インストール手順が最新のバージョンを反映しているか確認

### GitHub Release ノートを編集（オプション）

[Releases ページ](https://github.com/syeeel/koto-cli-go/releases) で最新リリースを編集：

- スクリーンショットやGIFを追加
- 主要な機能を強調
- 既知の問題があれば記載

### ソーシャルメディア（オプション）

- Twitter/X
- Reddit
- Hacker News
- Product Hunt

## 🔄 次回以降のリリース

次回からは、タグをプッシュするだけでOK：

```bash
# 変更をコミット
git add .
git commit -m "feat: Add awesome feature"
git push origin main

# タグを作成してプッシュ
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# あとは自動で処理されます！
```

## 🐛 トラブルシューティング

### ビルドが失敗した

```bash
# GitHub Actions のログを確認
# https://github.com/syeeel/koto-cli-go/actions

# よくある原因:
# - テストの失敗
# - go.mod/go.sum の不整合
# - .goreleaser.yaml の設定ミス
```

### Homebrew Formula が更新されない

```bash
# TAP_GITHUB_TOKEN が設定されているか確認
# https://github.com/syeeel/koto-cli-go/settings/secrets/actions

# トークンの権限を確認（public_repo が必要）

# homebrew-tap リポジトリが存在するか確認
# https://github.com/syeeel/homebrew-tap
```

### タグを間違えた

```bash
# ローカルのタグを削除
git tag -d v0.1.0

# リモートのタグを削除
git push origin :refs/tags/v0.1.0

# GitHubのReleasesページでリリースを削除

# 正しいタグを再作成
git tag -a v0.1.0 -m "..."
git push origin v0.1.0
```

## 📚 詳細ドキュメント

- [RELEASE.md](RELEASE.md) - 詳細なリリースガイド
- [SETUP_HOMEBREW.md](SETUP_HOMEBREW.md) - Homebrew Tap セットアップ
- [GoReleaser Documentation](https://goreleaser.com)

## ✅ チェックリスト

初回リリース前の最終確認：

- [ ] すべてのテストがパス
- [ ] ローカルビルドが成功
- [ ] `homebrew-tap` リポジトリを作成
- [ ] `TAP_GITHUB_TOKEN` を GitHub Secrets に追加
- [ ] CHANGELOG.md を更新
- [ ] README.md のインストール手順が正確
- [ ] コミット済み・プッシュ済み

リリース実行：

- [ ] タグを作成
- [ ] タグをプッシュ
- [ ] GitHub Actions の完了を確認
- [ ] GitHub Releases を確認
- [ ] Homebrew Tap を確認
- [ ] インストールをテスト

完了後：

- [ ] Release ノートを編集（オプション）
- [ ] アナウンス（オプション）

---

おめでとうございます！初回リリースが完了しました 🎉
