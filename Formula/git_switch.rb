class GitSwitch < Formula
  include Language::Python::Virtualenv
  desc "Description of your script"
  homepage "https://github.com/target-ops"
  version "v0.1.0-alpha"
  url "https://github.com/target-ops/GitSwitch/archive/refs/tags/v0.1.0-alpha.tar.gz"
  sha256 "5d1dcd2a20729b965fa7d0a20373dce6c9fa623e2a9e7a6a8dbe2e9cf8046b2e"
  depends_on "python@3.9"

  # resource "requests" do
  #   url "https://github.com/psf/requests/archive/refs/tags/v2.32.3.tar.gz"
  #   sha256 "f665576fc02d02e7b7f21630b915d70c14288f48decf76fad89b16a9f3975c42"
  # end
  
  def install
    bin.install "src"
    # system "cd #{prefix}/ && pip install -r requirements.txt"
    system "cd #{prefix}/bin/ && python3 -m src.main"
  end
end