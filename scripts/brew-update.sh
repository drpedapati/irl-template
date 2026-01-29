#!/bin/bash
# Update homebrew formula for a given version
# Usage: ./scripts/brew-update.sh 0.3.3

set -e

VERSION=$1
TAP_PATH="../homebrew-tap"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "Fetching SHA256 hashes for v$VERSION..."

ARM_SHA=$(gh release view v$VERSION --json assets -q '.assets[] | select(.name=="irl-darwin-arm64") | .digest' | sed 's/sha256://')
AMD_SHA=$(gh release view v$VERSION --json assets -q '.assets[] | select(.name=="irl-darwin-amd64") | .digest' | sed 's/sha256://')
LINUX_SHA=$(gh release view v$VERSION --json assets -q '.assets[] | select(.name=="irl-linux-amd64") | .digest' | sed 's/sha256://')

echo "  arm64:  $ARM_SHA"
echo "  amd64:  $AMD_SHA"
echo "  linux:  $LINUX_SHA"

cat > "$TAP_PATH/Formula/irl.rb" << EOF
class Irl < Formula
  desc "CLI for creating Idempotent Research Loop (IRL) projects"
  homepage "https://github.com/drpedapati/irl-template"
  version "$VERSION"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/drpedapati/irl-template/releases/download/v$VERSION/irl-darwin-arm64"
      sha256 "$ARM_SHA"
    else
      url "https://github.com/drpedapati/irl-template/releases/download/v$VERSION/irl-darwin-amd64"
      sha256 "$AMD_SHA"
    end
  end

  on_linux do
    url "https://github.com/drpedapati/irl-template/releases/download/v$VERSION/irl-linux-amd64"
    sha256 "$LINUX_SHA"
  end

  def install
    bin.install Dir["irl-*"].first => "irl"
  end

  test do
    system "#{bin}/irl", "--help"
  end
end
EOF

echo "Pushing formula update..."
cd "$TAP_PATH"
git add Formula/irl.rb
git commit -m "Update irl to v$VERSION"
git push

echo "✓ Formula updated to v$VERSION"
