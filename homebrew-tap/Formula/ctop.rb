# Homebrew formula for ctop - Top-like interface for container metrics
# https://github.com/eqms/ctop
#
# After a new release, update:
#   1. version variable
#   2. SHA256 checksums from the release's sha256sums.txt

class Ctop < Formula
  desc "Top-like interface for container metrics"
  homepage "https://github.com/eqms/ctop"
  version "0.8.7"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-arm64"
      sha256 "256da095cee6016dcc32a7caf935c2cb67f07ef0d165d704bd1aa762413c01ff"

      def install
        bin.install "ctop-#{version}-darwin-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-amd64"
      sha256 "aef1cad5f2df88be8eed2d5a51137634a58cddceb9b9a85e2dd4f62c48434fa7"

      def install
        bin.install "ctop-#{version}-darwin-amd64" => "ctop"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-arm64"
      sha256 "d788f093c1eac50d403d649bc4c2cd0c7f8c15da86e7acb5f59a7328a3e411b9"

      def install
        bin.install "ctop-#{version}-linux-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-amd64"
      sha256 "429457d9fc5ee1fd807f1675b5c2a3b4e8be8116d78bd1cde38480bfbfdbbf55"

      def install
        bin.install "ctop-#{version}-linux-amd64" => "ctop"
      end
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/ctop -v")
  end
end
