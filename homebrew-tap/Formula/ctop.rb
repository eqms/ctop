# Homebrew formula for ctop - Top-like interface for container metrics
# https://github.com/eqms/ctop
#
# After a new release, update:
#   1. version variable
#   2. SHA256 checksums from the release's sha256sums.txt

class Ctop < Formula
  desc "Top-like interface for container metrics"
  homepage "https://github.com/eqms/ctop"
  version "0.8.5"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-arm64"
      sha256 "b5da0eaff9cc45ca637f9feee1110dcab1f21c90e0860410fe1b76e49bb87579"

      def install
        bin.install "ctop-#{version}-darwin-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-darwin-amd64"
      sha256 "70507d2cc0c816467a3d228085633060bbe72252c3e4f4a00d3f2be35c4bb2e9"

      def install
        bin.install "ctop-#{version}-darwin-amd64" => "ctop"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-arm64"
      sha256 "3fd43a58b5108f111479e1b4d69c9635bf3036a713fd3aaaf2d3ff81aed2870c"

      def install
        bin.install "ctop-#{version}-linux-arm64" => "ctop"
      end
    else
      url "https://github.com/eqms/ctop/releases/download/v#{version}/ctop-#{version}-linux-amd64"
      sha256 "db310cca06ad5750eea80c8ecc0f88af559094e27976090b224e328ea0d8fdf4"

      def install
        bin.install "ctop-#{version}-linux-amd64" => "ctop"
      end
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/ctop -v")
  end
end
