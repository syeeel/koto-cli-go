# koto リリースガイド

このドキュメントでは、kotoの新しいバージョンをリリースする手順を説明します。

## 📋 リリース前チェックリスト

リリース前に、以下の項目を確認してください：

- [ ] すべてのテストがパス (`go test ./...`)
- [ ] コードがフォーマットされている (`go fmt ./...`)
- [ ] リントエラーがない (`go vet ./...`)
- [ ] CHANGELOG.mdが更新されている
- [ ] ドキュメントが最新
- [ ] ビルドが成功する (`go build -o bin/koto ./cmd/koto`)
- [ ] **GoReleaserのローカルテストが成功する** (`goreleaser release --snapshot --clean`)

## 🧪 ローカルでのテスト

本番リリース前に、GoReleaserの設定をローカルでテストできます：

```bash
# スナップショットビルド（GitHubにリリースせずにローカルでビルド）
goreleaser release --snapshot --clean

# 成功すると、dist/ディレクトリにバイナリが生成されます
ls -la dist/

# 各プラットフォーム向けのバイナリをテスト
./dist/koto_darwin_amd64_v1/koto --version
./dist/koto_linux_amd64_v1/koto --version
```

**設定ファイルの検証のみ:**

```bash
# .goreleaser.yamlの構文チェック
goreleaser check
```

**ドライラン（何も実行せず、実行内容だけ表示）:**

```bash
# リリース時に何が起こるかを確認
goreleaser release --skip=publish --clean
```

## 🚀 リリース手順

### 1. バージョンの決定

[セマンティックバージョニング](https://semver.org/)に従ってバージョンを決定します：

- **Major (X.0.0)**: 後方互換性のない変更
- **Minor (x.Y.0)**: 後方互換性のある機能追加
- **Patch (x.y.Z)**: 後方互換性のあるバグ修正

例: `v1.2.3`

### 2. CHANGELOGの更新

`CHANGELOG.md`に変更内容を記録します：

```markdown
## [1.2.3] - 2025-10-30

### Added
- 新機能の説明

### Changed
- 変更内容の説明

### Fixed
- 修正したバグの説明
```

### 3. タグの作成とプッシュ

```bash
# 新しいバージョンをタグ
git tag -a v1.0.0 -m "Release v1.0.0"

# タグをリモートにプッシュ
git push origin v1.0.0
```

### 4. 自動リリース

タグをプッシュすると、GitHub Actionsが自動的に以下を実行します：

1. **テストの実行** - すべてのテストを実行
2. **ビルド** - 各プラットフォーム向けにバイナリをビルド
   - macOS (Intel/Apple Silicon)
   - Linux (amd64/arm64)
   - Windows (amd64)
3. **アーカイブの作成** - tar.gz/zip形式でアーカイブ
4. **GitHub Releasesへの公開** - リリースページにバイナリとchangelogを公開
5. **チェックサムの生成** - セキュリティ検証用のチェックサムファイル作成

### 5. リリースの確認

GitHub Actionsのワークフローが完了したら、以下を確認します：

1. [Releases](https://github.com/syeeel/koto-cli-go/releases)ページを開く
2. 新しいリリースが作成されていることを確認
3. 各プラットフォーム向けのバイナリがアップロードされていることを確認
4. changelogが正しく表示されていることを確認

### 6. リリースノートの編集（オプション）

必要に応じて、GitHubのリリースページでリリースノートを編集します：

- 主要な変更点の強調
- 既知の問題の記載
- アップグレード手順（破壊的変更がある場合）
- スクリーンショットやGIFの追加

## 🔧 トラブルシューティング

### ビルドが失敗する

GitHub Actionsのログを確認します：

1. [Actions](https://github.com/syeeel/koto-cli-go/actions)ページを開く
2. 失敗したワークフローをクリック
3. エラーメッセージを確認

よくある原因：
- テストの失敗
- 依存関係の問題 (`go.mod`/`go.sum`の不整合)
- ビルドエラー

### タグを間違えた場合

ローカルとリモートのタグを削除します：

```bash
# ローカルのタグを削除
git tag -d v1.2.3

# リモートのタグを削除
git push origin :refs/tags/v1.2.3
```

その後、正しいタグを作成し直します。

### リリースを削除したい場合

1. GitHubのReleasesページでリリースを削除
2. タグも削除（上記参照）

## 📦 成果物

リリース後、以下が利用可能になります：

### GitHub Releases

すべてのプラットフォーム向けのバイナリ：

```
koto_v1.2.3_darwin_amd64.tar.gz
koto_v1.2.3_darwin_arm64.tar.gz
koto_v1.2.3_linux_amd64.tar.gz
koto_v1.2.3_linux_arm64.tar.gz
koto_v1.2.3_windows_amd64.zip
checksums.txt
```

### install.shスクリプト

ユーザーは以下のコマンドでインストールできます：

```bash
curl -sSfL https://raw.githubusercontent.com/syeeel/koto-cli-go/main/install.sh | sh
```

このスクリプトは自動的に最新リリースを取得してインストールします。

### Go install

Go環境があるユーザーは以下でインストールできます：

```bash
go install github.com/syeeel/koto-cli-go/cmd/koto@latest
```

## 🍺 Homebrew Tap のセットアップ

Homebrew Tapを使用すると、macOSユーザーが簡単にインストールできるようになります。

### 初回セットアップ手順

#### 1. homebrew-tapリポジトリを作成

GitHubで新しいパブリックリポジトリを作成します：

- リポジトリ名: **`homebrew-tap`** （必須の命名規則）
- オーナー: `syeeel`
- 公開設定: Public
- 初期化: README を追加（オプション）

```bash
# CLIで作成する場合（gh コマンドが必要）
gh repo create syeeel/homebrew-tap --public --description "Homebrew tap for koto"
```

#### 2. Personal Access Token を作成

1. GitHub で [Settings → Developer settings → Personal access tokens → Tokens (classic)](https://github.com/settings/tokens)
2. "Generate new token (classic)" をクリック
3. 設定:
   - **Note**: `GoReleaser - koto homebrew-tap`
   - **Expiration**: 好みに応じて（推奨: No expiration または 1 year）
   - **Select scopes**:
     - ✅ `public_repo` （パブリックリポジトリの場合）
     - または ✅ `repo` （プライベートリポジトリの場合）
4. "Generate token" をクリック
5. **トークンをコピー**（このページを離れると二度と表示されません）

#### 3. GitHub シークレットに追加

1. [koto-cli-go リポジトリの Settings → Secrets and variables → Actions](https://github.com/syeeel/koto-cli-go/settings/secrets/actions)
2. "New repository secret" をクリック
3. 設定:
   - **Name**: `TAP_GITHUB_TOKEN`
   - **Secret**: 上記でコピーしたトークン
4. "Add secret" をクリック

#### 4. 設定ファイルの確認

以下のファイルはすでに設定済みです：

- ✅ `.goreleaser.yaml` - Homebrew Tap設定が有効
- ✅ `.github/workflows/release.yml` - TAP_GITHUB_TOKEN を使用

#### 5. 動作確認

次回のリリース時（タグをプッシュ）に、GoReleaserが自動的に：

1. `homebrew-tap` リポジトリに Formula を作成
2. `Formula/koto.rb` ファイルを更新
3. バージョン情報とダウンロードURLを自動更新

#### 6. ユーザー向けインストール手順

セットアップ完了後、ユーザーは以下でインストールできます：

```bash
brew tap syeeel/tap
brew install koto
```

### トラブルシューティング

**エラー: `401 Bad credentials`**

- TAP_GITHUB_TOKEN が正しく設定されているか確認
- トークンに `public_repo` または `repo` スコープがあるか確認
- トークンの有効期限が切れていないか確認

**エラー: `404 Not Found`**

- `homebrew-tap` リポジトリが存在するか確認
- リポジトリが Public になっているか確認

**Formula が更新されない**

- GitHub Actions のログを確認
- `.goreleaser.yaml` の `repository.owner` と `repository.name` が正しいか確認

## 📝 リリーススケジュール（例）

プロジェクトの規模に応じて、リリーススケジュールを決めることができます：

- **Major**: 年1-2回（大きな機能追加や破壊的変更）
- **Minor**: 月1-2回（新機能追加）
- **Patch**: 随時（バグ修正）

## 🔗 参考リンク

- [GoReleaser Documentation](https://goreleaser.com)
- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Homebrew Tap Guide](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)

---

リリースプロセスに問題がある場合は、[Issues](https://github.com/syeeel/koto-cli-go/issues)で報告してください。
