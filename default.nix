{ pkgs, version ? "dev" }:

pkgs.buildGoModule rec {
  pname = "krampus";
  inherit version;

  src = ./.;

  vendorHash = "sha256-3WfBvlM4aCwP9YqUWU4lnxpSLWHQNGWXjjg2F91awnY=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
