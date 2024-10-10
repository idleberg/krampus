{
  description = "Command-line tool to kill one or more processes by its port number.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; }; in
      rec {
        packages.default = import ./default.nix { inherit pkgs; };
      }) // {
        overlays.default = final: prev: {
          krampus = import ./default.nix { pkgs = final; };
        };
      };
}
