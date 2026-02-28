# Homebrew formula for ctop - Top-like interface for container metrics
# https://github.com/eqms/ctop
#
# After a new release, update:
#   1. version variable
#   2. SHA256 checksums from the release's sha256sums.txt

class Ctop < Formula
  desc "Top-like interface for container metrics"
  homepage "https://github.com/eqms/ctop"
  version "0.8.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-arm64"
      sha256 "e85f1dcd01dd02a30b128510169d1a95f4e9902db935d61365e44633113e5d63"

      def install
        bin.install "ctop-#{version}-darwin-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-amd64"
      sha256 "4d74484a534c3c2fab38d265a1d7c666accb84117e5f5d2f899392307befb6f6"

      def install
        bin.install "ctop-#{version}-darwin-amd64" => "ctop"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-arm64"
      sha256 "d1d661260f67e6abff8468edf2df9409f4382768e11de3a2c3e61d57cfa7e750"

      def install
        bin.install "ctop-#{version}-linux-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-amd64"
      sha256 "aa48e329a872e4e59f8b287b5bdb692650a56695dec6c293649db1ed6a211d3d"

      def install
        bin.install "ctop-#{version}-linux-amd64" => "ctop"
      end
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/ctop -v")
  end
end
