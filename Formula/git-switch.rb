class GitSwitch < Formula
  include Language::Python::Virtualenv

  desc "Description of your script"
  homepage "https://github.com/target-ops"
  url "https://github.com/target-ops/GitSwitch/archive/refs/tags/v0.1.0-alpha.zip"
  sha256 "e534e60aea8c8465911e03659c97d56aab4330b6a86b85e408f5fc1e4d8c4250"
  license "LICENSE_TYPE"

  depends_on "python@3.9"

  def install
    virtualenv_install_with_resources
  end
  # test do
  #   system "#{bin}/your_script_name", "--version"
  # end
end