class Retrier < Formula
    desc "A tool for retrying commands until they succeed"
    homepage "https://github.com/dantheman0207/retrier"
    url "https://github.com/dantheman0207/homebrew-retrier/releases/download/v0.1.5/retrier-v0.1.5-darwin-arm.tar.gz"
    sha256 "ee41319e62db5c64a1bbfdff627287d01e227f8bcff6178a47dbba46f5c8bdea"
    license "MIT"
  
    depends_on "go" => :build
  
    def install
      system "go", "build", *std_go_args
    end
  
    test do
      # Basic test to verify the tool runs
      assert_match "Usage", shell_output("#{bin}/retrier --help")
    end
  end