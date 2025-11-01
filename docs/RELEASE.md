# koto Release Guide

This document describes the procedure for releasing new versions of koto.

## 📋 Pre-Release Checklist

Before releasing, verify the following items:

- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No lint errors (`go vet ./...`)
- [ ] CHANGELOG.md is updated
- [ ] Documentation is up to date
- [ ] Build succeeds (`go build -o bin/koto ./cmd/koto`)
- [ ] **GoReleaser local test succeeds** (`goreleaser release --snapshot --clean`)

## 🧪 Local Testing

Before production release, you can test GoReleaser configuration locally:

```bash
# Snapshot build (build locally without releasing to GitHub)
goreleaser release --snapshot --clean

# On success, binaries are generated in the dist/ directory
ls -la dist/

# Test binaries for each platform
./dist/koto_darwin_amd64_v1/koto --version
./dist/koto_linux_amd64_v1/koto --version
```

**Configuration file validation only:**

```bash
# Check .goreleaser.yaml syntax
goreleaser check
```

**Dry run (show what would be executed without actually executing):**

```bash
# See what happens during release
goreleaser release --skip=publish --clean
```

## 🚀 Release Procedure

### 1. Determine Version

Decide the version following [Semantic Versioning](https://semver.org/):

- **Major (X.0.0)**: Breaking changes
- **Minor (x.Y.0)**: Backward-compatible feature additions
- **Patch (x.y.Z)**: Backward-compatible bug fixes

Example: `v1.2.3`

### 2. Update CHANGELOG

Record changes in `CHANGELOG.md`:

```markdown
## [1.0.0] - 2025-10-30

### Added
- Description of new features

### Changed
- Description of changes

### Fixed
- Description of bug fixes
```

### 3. Create and Push Tag

```bash
# Tag the new version
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag to remote
git push origin v1.0.0
```

### 4. Automatic Release

When you push a tag, GitHub Actions automatically executes the following:

1. **Run Tests** - Execute all tests
2. **Build** - Build binaries for each platform
   - macOS (Intel/Apple Silicon)
   - Linux (amd64/arm64)
   - Windows (amd64)
3. **Create Archives** - Create archives in tar.gz/zip format
4. **Publish to GitHub Releases** - Publish binaries and changelog to releases page
5. **Generate Checksums** - Create checksum files for security verification

### 5. Verify Release

After GitHub Actions workflow completes, verify the following:

1. Open the [Releases](https://github.com/syeeel/koto-cli-go/releases) page
2. Confirm new release is created
3. Confirm binaries for each platform are uploaded
4. Confirm changelog is displayed correctly

### 6. Edit Release Notes (Optional)

If necessary, edit release notes on GitHub's release page:

- Highlight major changes
- Document known issues
- Upgrade instructions (if there are breaking changes)
- Add screenshots or GIFs

## 🔧 Troubleshooting

### Build Fails

Check GitHub Actions logs:

1. Open the [Actions](https://github.com/syeeel/koto-cli-go/actions) page
2. Click on the failed workflow
3. Check error messages

Common causes:
- Test failures
- Dependency issues (`go.mod`/`go.sum` inconsistencies)
- Build errors

### Wrong Tag

Delete local and remote tags:

```bash
# Delete local tag
git tag -d v1.2.3

# Delete remote tag
git push origin :refs/tags/v1.2.3
```

Then recreate the correct tag.

### Delete Release

1. Delete the release on GitHub's Releases page
2. Also delete the tag (see above)

## 📦 Artifacts

After release, the following become available:

### GitHub Releases

Binaries for all platforms:

```
koto_v1.2.3_darwin_amd64.tar.gz
koto_v1.2.3_darwin_arm64.tar.gz
koto_v1.2.3_linux_amd64.tar.gz
koto_v1.2.3_linux_arm64.tar.gz
koto_v1.2.3_windows_amd64.zip
checksums.txt
```

### install.sh Script

Users can install with the following command:

```bash
curl -sSfL https://raw.githubusercontent.com/syeeel/koto-cli-go/main/install.sh | sh
```

This script automatically fetches and installs the latest release.

### Go install

Users with Go environment can install with:

```bash
go install github.com/syeeel/koto-cli-go/cmd/koto@latest
```

## 🍺 Homebrew Tap Setup

Using Homebrew Tap makes it easy for macOS users to install.

### Initial Setup Procedure

#### 1. Create homebrew-tap Repository

Create a new public repository on GitHub:

- Repository name: **`homebrew-tap`** (required naming convention)
- Owner: `syeeel`
- Visibility: Public
- Initialize: Add README (optional)

```bash
# Create via CLI (requires gh command)
gh repo create syeeel/homebrew-tap --public --description "Homebrew tap for koto"
```

#### 2. Create Personal Access Token

1. On GitHub, go to [Settings → Developer settings → Personal access tokens → Tokens (classic)](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Settings:
   - **Note**: `GoReleaser - koto homebrew-tap`
   - **Expiration**: As preferred (recommended: No expiration or 1 year)
   - **Select scopes**:
     - ✅ `public_repo` (for public repositories)
     - or ✅ `repo` (for private repositories)
4. Click "Generate token"
5. **Copy the token** (it won't be shown again once you leave this page)

#### 3. Add to GitHub Secrets

1. Go to [koto-cli-go repository Settings → Secrets and variables → Actions](https://github.com/syeeel/koto-cli-go/settings/secrets/actions)
2. Click "New repository secret"
3. Settings:
   - **Name**: `TAP_GITHUB_TOKEN`
   - **Secret**: The token copied above
4. Click "Add secret"

#### 4. Verify Configuration Files

The following files are already configured:

- ✅ `.goreleaser.yaml` - Homebrew Tap configuration enabled
- ✅ `.github/workflows/release.yml` - Uses TAP_GITHUB_TOKEN

#### 5. Verify Operation

On the next release (push tag), GoReleaser will automatically:

1. Create Formula in `homebrew-tap` repository
2. Update `Formula/koto.rb` file
3. Automatically update version information and download URL

#### 6. User Installation Instructions

After setup completion, users can install with:

```bash
brew tap syeeel/tap
brew install koto
```

### Troubleshooting

**Error: `401 Bad credentials`**

- Verify TAP_GITHUB_TOKEN is correctly set
- Verify token has `public_repo` or `repo` scope
- Verify token has not expired

**Error: `404 Not Found`**

- Verify `homebrew-tap` repository exists
- Verify repository is Public

**Formula Not Updated**

- Check GitHub Actions logs
- Verify `repository.owner` and `repository.name` in `.goreleaser.yaml` are correct

## 📝 Release Schedule (Example)

Depending on project scale, you can decide on a release schedule:

- **Major**: 1-2 times per year (major features or breaking changes)
- **Minor**: 1-2 times per month (new features)
- **Patch**: As needed (bug fixes)

## 🔗 References

- [GoReleaser Documentation](https://goreleaser.com)
- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Homebrew Tap Guide](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)

---

If you have any issues with the release process, please report them in [Issues](https://github.com/syeeel/koto-cli-go/issues).
