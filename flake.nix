{
  inputs = {
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, utils }: utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShell = pkgs.mkShell {
        TF_CLI_CONFIG_FILE = "$$PWD/.terraformrc";

        buildInputs = with pkgs; [
          go
          gnumake
          python311Packages.keyring
          opentofu
          goreleaser
        ];
      };
    }
  );
}
