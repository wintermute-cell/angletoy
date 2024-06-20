{ pkgs ? import (fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/nixos-23.11.tar.gz";
  }) {}
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    (buildGoModule rec {
      pname = "longgopher";
      version = "0.0.3";
      src = fetchFromGitHub {
        owner = "sheepla";
        repo = "longgopher";
        rev = "v${version}";
        sha256 = "sha256-q0I53lIrJBEx171iEDzrkemjQbDvhi0K8Snhwso4K5Y=";
      };
      vendorHash = "sha256-nzPHx+c369T4h9KETqMurxZK3LsJAhwBaunkcWIW3Ps=";
      subPackages = [ "." ];
    })
    libGL
    xorg.libXi
    xorg.libXcursor
    xorg.libXrandr
    xorg.libXinerama
    wayland
    libxkbcommon
  ];

  # this is needed for delve to work with cgo
  # see: https://wiki.nixos.org/wiki/Go#Using_cgo_on_NixOS
  hardeningDisable = [ "fortify" ];

  shellHook = ''
    if [ -z "$IN_DEV_SHELL" ]; then
      echo -e "\033[1;32mEntering Nix shell...\033[0m"
      export IN_DEV_SHELL=1
      export PS1="[Nix] $PS1"
      longgopher -l 5
      echo ""
    else
      echo -e "\033[1;31mAlready in Nix shell!\033[0m"
      exit 1
    fi
  '';
}

