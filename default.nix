{pkgs}:
pkgs.buildGoModule rec {
  pname = "krampus";
  version = "0.3.0";

  src = ./.;

  vendorHash = "sha256-3WfBvlM4aCwP9YqUWU4lnxpSLWHQNGWXjjg2F91awnY=";

  ldflags = ["-s" "-w" "-X=main.Version=${version}"];
}
