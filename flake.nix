{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

  outputs =
    {
      self,
      nixpkgs,
      ...
    }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f (import nixpkgs { inherit system; }));
    in
    {

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          name = "real";
          buildInputs = with pkgs; [
            gnumake
            go
            gofumpt
            golines
            gopls
            revive
          ];
        };

      });
    };
}
