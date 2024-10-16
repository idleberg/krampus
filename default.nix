{ pkgs }:

pkgs.buildGoModule rec {
  pname = "krampus";
  version = "0.2.1";

  src = ./.;

  vendorHash = "sha256-3WfBvlM4aCwP9YqUWU4lnxpSLWHQNGWXjjg2F91awnY=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
