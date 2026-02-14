# 013: tsk Homebrew tap

## Objective
Create a Homebrew tap so users can install tsk via `brew install zarldev/tap/tsk`.

## Context
Spec 012 adds a release pipeline that produces cross-platform archives on GitHub Releases. This spec adds a Homebrew tap that points to those release assets, and updates the release workflow to auto-update the formula on each release.

## Requirements

### New repo: `zarldev/homebrew-tap`
Create a Homebrew tap repository with a single formula.

### Formula: `Formula/tsk.rb`
```ruby
class Tsk < Formula
  desc "Simple CLI task tracker"
  homepage "https://zarldev.github.io/tsk/"
  version "VERSION"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/zarldev/tsk/releases/download/vVERSION/tsk-darwin-arm64.tar.gz"
      sha256 "SHA256"
    else
      url "https://github.com/zarldev/tsk/releases/download/vVERSION/tsk-darwin-amd64.tar.gz"
      sha256 "SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/zarldev/tsk/releases/download/vVERSION/tsk-linux-arm64.tar.gz"
      sha256 "SHA256"
    else
      url "https://github.com/zarldev/tsk/releases/download/vVERSION/tsk-linux-amd64.tar.gz"
      sha256 "SHA256"
    end
  end

  def install
    bin.install "tsk"
  end

  test do
    assert_match "tsk", shell_output("#{bin}/tsk version")
  end
end
```

VERSION and SHA256 are placeholders — the release workflow will substitute them.

### Update release workflow in `zarldev/tsk`
Add a step at the end of `.github/workflows/release.yml` (from spec 012) that:
1. Computes sha256 checksums for all archives
2. Clones `zarldev/homebrew-tap`
3. Updates `Formula/tsk.rb` with the new version and checksums
4. Commits and pushes the updated formula

This step needs write access to the tap repo. Use `GITHUB_TOKEN` if the tap is in the same org, or a PAT if needed. Since both repos are under `zarldev`, the default `GITHUB_TOKEN` with `contents: write` on the tap repo should work — set up the tap repo to allow GitHub Actions from `zarldev/tsk`.

Alternatively, use a repository dispatch event: the release workflow sends a dispatch to `homebrew-tap` with the version and checksums, and a workflow in the tap repo updates the formula. This is cleaner but more complex. Pick whichever is simpler.

### README for homebrew-tap
Simple README:
```markdown
# homebrew-tap

Homebrew formulae for zarldev tools.

## Install

    brew install zarldev/tap/tsk

## Update

    brew upgrade tsk
```

### Update tsk docs
Update `docs/content/docs.md` install section to include Homebrew:
```
## install

### homebrew (macOS/Linux)

    brew install zarldev/tap/tsk

### go install

    go install github.com/zarldev/tsk/cmd/tsk@latest
```

Also update `README.md` in zarldev/tsk with the brew install option.

### Update tsk landing page
Update `docs/layouts/index.html` — change the install terminal window to show both options or just the brew option (simpler for most users).

## Target Repo
zarldev/tsk (release workflow update + docs) and zarldev/homebrew-tap (new repo)

## Agent Role
backend

## Files to Create (in zarldev/homebrew-tap)
- `Formula/tsk.rb`
- `README.md`

## Files to Modify (in zarldev/tsk)
- `.github/workflows/release.yml` (add tap update step)
- `docs/content/docs.md` (add brew install)
- `README.md` (add brew install)
- `docs/layouts/index.html` (update install block)

## Dependencies
- 012 (release pipeline) must be merged

## Notes
- The formula will be a template initially. First real release (from 012) populates the version and checksums.
- Keep it simple — the tap repo is just a formula and a README.
- Test with `brew install --build-from-source` locally if possible.
- No license file exists yet in tsk — consider adding MIT LICENSE as part of this or a separate item.
