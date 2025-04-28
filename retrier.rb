class Retrier < Formula
    desc "A tool for retrying commands until they succeed"
    homepage "https://github.com/dantheman0207/retrier"
    url "https://github.com/dantheman0207/homebrew-retrier/releases/download/v0.1.14/retrier-v0.1.14-darwin-arm.tar.gz"
    sha256 "5a9d88627f236f5f7ed0437d89e0d40b7a57ba9a902fd654f4b7b3ac7ba0b4a4"
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