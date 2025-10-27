{
  description = "A simple flake for building ezazssh, a Go-based CLI tool.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        packages.ezazssh = pkgs.buildGoModule {
          pname = "ezazssh";
          version = "v${self.shortRev or "dev"}";
          env.CGO_ENABLED = 0;
          ldflags = ["-w" "-s" "-extldflags \"-static\"" ];
          src = self;
#          src = pkgs.fetchFromGitHub {
#            owner = "burneystarke";
#            repo = "ezazssh";
#            rev = self.rev or "main";
#            # You can run `nix flake update` to refresh the hash automatically
#            sha256 = "sha256-d3IgW3eBlIR0+fSoBJljD3eZQg2VVP8OPfAibCjGhS4=";
#          };
          vendorHash = "sha256-NTLXGC61axh/8MtfOIpKsk+ykTz+YkrLX4UCjddBlYw=";
          subPackages = [ "." ];

          meta = with pkgs.lib; {
            description = "Azure SSH utility written in Go";
            homepage = "https://github.com/burneystarke/ezazssh";
            license = licenses.mit;
            maintainers = [ maintainers.yourGithubUsername ];
            platforms = platforms.all;
          };
        };

        defaultPackage = self.packages.${system}.ezazssh;
        defaultApp = flake-utils.lib.mkApp {
          drv = self.packages.${system}.ezazssh;
        };
      });
}
