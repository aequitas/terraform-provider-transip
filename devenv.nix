{ pkgs, lib,... }:

{       
  languages.go.enable = true;

  packages = [ 
    pkgs.gnumake
    pkgs.python311Packages.keyring 
    pkgs.terraform
    pkgs.opentofu
  ];

  env.TF_CLI_CONFIG_FILE = "$$PWD/.terraformrc";
}
