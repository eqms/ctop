# Homebrew formula for ctop - Top-like interface for container metrics
# https://github.com/eqms/ctop
#
# After a new release, update:
#   1. version variable
#   2. SHA256 checksums from the release's sha256sums.txt

class Ctop < Formula
  desc "Top-like interface for container metrics"
  homepage "https://github.com/eqms/ctop"
  version "0.8.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-arm64"
      sha256 "PLACEHOLDER_DARWIN_ARM64"

      def install
        bin.install "ctop-#{version}-darwin-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-amd64"
      sha256 "PLACEHOLDER_DARWIN_AMD64"

      def install
        bin.install "ctop-#{version}-darwin-amd64" => "ctop"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-arm64"
      sha256 "PLACEHOLDER_LINUX_ARM64"

      def install
        bin.install "ctop-#{version}-linux-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-amd64"
      sha256 "PLACEHOLDER_LINUX_AMD64"

      def install
        bin.install "ctop-#{version}-linux-amd64" => "ctop"
      end
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/ctop -v")
  end
end
