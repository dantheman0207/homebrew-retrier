class Retrier < Formula
    desc "A tool for retrying commands until they succeed"
    homepage "https://github.com/dantheman0207/retrier"
    url "https://github.com/dantheman0207/homebrew-retrier/releases/download/v0.1.18/retrier-v0.1.18-darwin-arm.tar.gz"
    sha256 "d940e6704c26cb5b44fe3cabbf6de6b73d73bebeb9fe4c425a0eda48ad06f88d"
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