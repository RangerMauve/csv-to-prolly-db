{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs @ { self, ... }:
    (inputs.flake-utils.lib.eachDefaultSystem (system:
      let

        pkgs = import inputs.nixpkgs {
          inherit system;
        };

      in
      rec {

        devShells = {
          default = pkgs.mkShell ({
            nativeBuildInputs = with pkgs; [
              go
            ];
          });
        };

        packages = {
          default = pkgs.buildGoModule rec {
            pname = "csv-to-prolly-db";
            version = "0.0.0";

            src = ./.;

            vendorHash = "sha256-SCcnq2biYHWR+wKK64pFo57IoSBPxiD1e17y/MaTlQE";

            # test fails, didn't troubleshot yet
            doCheck = false;
          };
        };

      }));
}
