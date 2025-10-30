# Homebrew Tap セットアップガイド

このガイドでは、koto用のHomebrew Tapをセットアップし、macOSユーザーが `brew install` でインストールできるようにする手順を説明します。

## 📋 前提条件

- GitHubアカウント
- koto-cli-goリポジトリへの管理者権限
- GoReleaserがセットアップ済み（`.goreleaser.yaml`が設定済み）

## 🚀 セットアップ手順

### ステップ 1: homebrew-tap リポジトリの作成

#### 方法A: GitHub CLI を使用（推奨）

```bash
gh repo create syeeel/homebrew-tap --public --description "Homebrew tap for koto CLI"
```

#### 方法B: GitHubウェブUI を使用

1. [GitHub](https://github.com/new) で新しいリポジトリを作成
2. 設定:
   - **Repository name**: `homebrew-tap` ⚠️ 必ず `homebrew-` で始める必要があります
   - **Owner**: `syeeel`
   - **Public** を選択
   - **Initialize this repository with**:
     - ✅ Add a README file（オプション）
     - ✅ Add .gitignore: None
     - ✅ Choose a license: MIT（オプション）
3. "Create repository" をクリック

### ステップ 2: Personal Access Token (PAT) の作成

#### 2.1 トークンの生成

1. GitHub の [Personal Access Tokens (Classic)](https://github.com/settings/tokens) ページを開く
2. "Generate new token" → "Generate new token (classic)" をクリック
3. トークン設定:

   | 項目 | 設定値 |
   |------|--------|
   | **Note** | `GoReleaser - koto homebrew-tap` |
   | **Expiration** | `No expiration` または `1 year` |
   | **Scopes** | `public_repo` ✅ |

4. "Generate token" をクリック
5. **重要**: 表示されたトークンをコピー（このページを離れると二度と表示されません）

#### 2.2 トークンの保存（オプションだが推奨）

安全な場所にトークンを保存しておくことを推奨します：

```bash
# パスワードマネージャーに保存
# または、ローカルに暗号化して保存（例）
echo "ghp_xxxxxxxxxxxx" > ~/.ssh/github_tap_token
chmod 600 ~/.ssh/github_tap_token
```

### ステップ 3: GitHub Secrets への追加

#### 3.1 koto-cli-go リポジトリの Secrets 設定

1. [koto-cli-go リポジトリ](https://github.com/syeeel/koto-cli-go) を開く
2. **Settings** → **Secrets and variables** → **Actions** をクリック
3. "New repository secret" をクリック
4. シークレット設定:

   | 項目 | 設定値 |
   |------|--------|
   | **Name** | `TAP_GITHUB_TOKEN` |
   | **Secret** | ステップ2でコピーしたトークン |

5. "Add secret" をクリック

#### 3.2 設定の確認

Secrets が正しく追加されたか確認：

- Settings → Secrets and variables → Actions
- `TAP_GITHUB_TOKEN` が表示されているはずです（値は隠されています）

### ステップ 4: 設定ファイルの確認

以下のファイルは既に設定済みです（確認のみ）：

#### `.goreleaser.yaml`

```yaml
brews:
  - repository:
      owner: syeeel
      name: homebrew-tap
      token: "{{ .Env.TAP_GITHUB_TOKEN }}"
    directory: Formula
    homepage: "https://github.com/syeeel/koto-cli-go"
    description: "Interactive ToDo management CLI tool with beautiful TUI"
    license: "MIT"
```

#### `.github/workflows/release.yml`

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  TAP_GITHUB_TOKEN: ${{ secrets.TAP_GITHUB_TOKEN }}
```

### ステップ 5: 初回リリース

#### 5.1 リリースタグの作成

```bash
# 最新の状態を確認
git status
git pull origin main

# 初回リリースタグを作成
git tag -a v0.1.0 -m "Initial release with Homebrew support"

# タグをプッシュ
git push origin v0.1.0
```

#### 5.2 GitHub Actions の確認

1. [Actions タブ](https://github.com/syeeel/koto-cli-go/actions) を開く
2. "Release" ワークフローが実行されているか確認
3. ワークフローをクリックして詳細を確認

ワークフローが成功すると：
- ✅ GitHub Releases にバイナリが公開
- ✅ `homebrew-tap` リポジトリに `Formula/koto.rb` が作成
- ✅ Homebrew Formula が自動更新

#### 5.3 homebrew-tap リポジトリの確認

1. [homebrew-tap リポジトリ](https://github.com/syeeel/homebrew-tap) を開く
2. `Formula/koto.rb` ファイルが作成されているか確認
3. ファイルの内容を確認（バージョン、URL、SHA256ハッシュが自動生成されている）

### ステップ 6: インストールのテスト

#### 6.1 Tap の追加

```bash
brew tap syeeel/tap
```

#### 6.2 koto のインストール

```bash
brew install koto
```

#### 6.3 動作確認

```bash
# バージョン確認
koto --version

# アプリ起動
koto
```

## 🎉 完了！

これで、Homebrew Tap のセットアップが完了しました。

今後のリリースでは、タグをプッシュするだけで自動的にHomebrew Formulaが更新されます。

## 📚 ユーザー向けインストール手順

セットアップ完了後、ユーザーには以下の手順を案内できます：

```bash
# Tap を追加
brew tap syeeel/tap

# koto をインストール
brew install koto

# アップデート
brew upgrade koto

# アンインストール
brew uninstall koto
```

## 🔧 トラブルシューティング

### エラー: `401 Bad credentials`

**原因**: GitHub トークンが無効または権限不足

**解決策**:
1. トークンの有効期限を確認
2. トークンに `public_repo` スコープがあるか確認
3. GitHub Secrets が正しく設定されているか確認
4. 必要に応じて新しいトークンを生成して再設定

### エラー: `404 Not Found`

**原因**: `homebrew-tap` リポジトリが見つからない

**解決策**:
1. リポジトリが存在するか確認: https://github.com/syeeel/homebrew-tap
2. リポジトリが **Public** になっているか確認
3. リポジトリ名が `homebrew-tap` であるか確認（`homebrew-` プレフィックスが必須）

### Formula が更新されない

**原因**: GoReleaser の設定ミスまたはワークフローエラー

**解決策**:
1. GitHub Actions のログを確認
2. `.goreleaser.yaml` の `repository.owner` と `repository.name` が正しいか確認
3. ローカルでテスト: `goreleaser release --snapshot --clean`

### Homebrew でインストールできない

**原因**: Formula の構文エラーまたはバイナリURL問題

**解決策**:
```bash
# Formula の検証
brew audit --strict syeeel/tap/koto

# Formula の詳細確認
brew info syeeel/tap/koto

# キャッシュのクリア
brew cleanup
rm -rf ~/Library/Caches/Homebrew/koto*

# 再インストール
brew reinstall syeeel/tap/koto
```

## 🔄 更新フロー

新しいバージョンをリリースする際の流れ：

```bash
# 1. 変更をコミット
git add .
git commit -m "feat: Add new feature"

# 2. タグを作成
git tag -a v0.2.0 -m "Release v0.2.0"

# 3. タグをプッシュ（これだけでOK！）
git push origin v0.2.0

# 4. GitHub Actions が自動的に：
#    - ビルド
#    - GitHub Releases に公開
#    - Homebrew Formula を更新
```

ユーザー側では：

```bash
# 最新版に更新
brew update
brew upgrade koto
```

## 📖 参考リンク

- [Homebrew Tap 公式ドキュメント](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [GoReleaser Homebrew 設定](https://goreleaser.com/customization/homebrew/)
- [GitHub Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

## ❓ FAQ

### Q: 複数のプロジェクトで同じ Tap を使えますか？

A: はい、1つの `homebrew-tap` リポジトリに複数の Formula を追加できます。

### Q: プライベートリポジトリでも使えますか？

A: はい、ただし PAT のスコープを `repo`（full control）に変更する必要があります。

### Q: Tap の名前を変更できますか？

A: 変更可能ですが、以下の制約があります：
- 必ず `homebrew-` で始まる必要があります
- ユーザーは `brew tap <username>/<tap-name>` で Tap を追加します
  - 例: `homebrew-koto` → `brew tap syeeel/koto`

### Q: Cask（GUIアプリ）もサポートできますか？

A: kotoはCLIツールなので Formula を使用していますが、GUIアプリの場合は Cask を使用します。

---

セットアップに問題がある場合は、[Issues](https://github.com/syeeel/koto-cli-go/issues) で質問してください。
