# Homebrew formula for ctop - Top-like interface for container metrics
# https://github.com/eqms/ctop
#
# After a new release, update:
#   1. version variable
#   2. SHA256 checksums from the release's sha256sums.txt

class Ctop < Formula
  desc "Top-like interface for container metrics"
  homepage "https://github.com/eqms/ctop"
  version "0.8.4"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-arm64"
      sha256 "2ef960b7fd1b550680c67fe5cbc064a384229d12890cef21b0b1faeed714266f"

      def install
        bin.install "ctop-#{version}-darwin-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-amd64"
      sha256 "9d5b0b7d439062f5bcfff49628dde3a871e7d56d6c9779e42f5ed0ecc0e11e9f"

      def install
        bin.install "ctop-#{version}-darwin-amd64" => "ctop"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-arm64"
      sha256 "448f33f3e0c4d699997d1c78260cb3b97df3d88f7d0f60bd23ea1fa810376ce7"

      def install
        bin.install "ctop-#{version}-linux-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-amd64"
      sha256 "e39d2e6baf0f5d77058d6fc5f85d71db9a7e3ab88814c0009e97cd8a5a163e09"

      def install
        bin.install "ctop-#{version}-linux-amd64" => "ctop"
      end
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/ctop -v")
  end
end
