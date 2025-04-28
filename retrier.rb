class Retrier < Formula
    desc "A tool for retrying commands until they succeed"
    homepage "https://github.com/dantheman0207/retrier"
    url "https://github.com/dantheman0207/homebrew-retrier/releases/download/v0.1.13/retrier-v0.1.13-darwin-arm.tar.gz"
    sha256 "522dc1f861652b0080b2add4df08d4feb3ba6297728f24a597f0bc289a1cae2e"
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