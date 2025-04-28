class Retrier < Formula
    desc "A tool for retrying commands until they succeed"
    homepage "https://github.com/dantheman0207/retrier"
    url "https://github.com/dantheman0207/homebrew-retrier/releases/download/v0.1.10/retrier-v0.1.10-darwin-arm.tar.gz"
    sha256 "2379f68cf5a88b3dc3e773cbe33cdfc91a7a4c163b4816e8af969f2d87beff53"
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