{ writeShellScriptBin
, ainglenixPath
}:

{
  ciSetupNixConf = writeShellScriptBin "hc-ci-setup-nix-conf.sh" ''
    ${ainglenixPath}/ci/setup-hydra-cache.sh
  '';
}
