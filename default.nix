{ pkgs }:

pkgs.buildGoModule rec {
  pname = "krampus";
  version = "0.2.0";

  src = ./.;

  vendorHash = "sha256-SzCpleSRUYSBHuvJA7Rs49IvTjI5ILBW68QkYYCPqZ4=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
