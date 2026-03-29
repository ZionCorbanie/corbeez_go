{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gnumake
    cargo
    rustc
    tailwindcss_4
    tailwindcss
    nodejs_24
  ];

  shellHook = ''
  '';
}

