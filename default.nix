{ pkgs }:

pkgs.buildGoModule rec {
  pname = "krampus";
  version = "0.2.1";

  src = ./.;

  vendorHash = "sha256-c6aTTAKfk0h2r51wnDicKNi7iilT4SGjNxHawAD1GYY=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
