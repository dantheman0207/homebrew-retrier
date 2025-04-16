class Retrier < Formula
    desc "A tool for retrying commands until they succeed"
    homepage "https://github.com/dantheman0207/retrier"
    url "https://github.com/dantheman0207/retrier/archive/v0.1.4.tar.gz"
    sha256 "0ddf6de945c438113a754b436b7ada30aea1581f179e2006fedaa4e75ce0879c"
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