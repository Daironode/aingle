# aingle_cli

Provides the `ai` binary, a helpful CLI tool for working with AIngle.

```shell
$ ai -h

aingle_cli 0.1.0
AIngle CLI

Work with SAF and hApp bundle files, set up sandbox environments for testing and development purposes, make direct admin
calls to running conductors, and more.

USAGE:
    ai <SUBCOMMAND>

FLAGS:
    -h, --help       Prints help information
    -V, --version    Prints version information

SUBCOMMANDS:
    app        Work with hApp bundles
    saf        Work with SAF bundles
    help       Prints this message or the help of the given subcommand(s)
    sandbox    Work with sandboxed environments for testing and development
```

## Docs

Each top-level subcommand is implemented as a separate crate. See:

- [aingle_cli_bundle](https://github.com/AIngleLab/aingle/tree/develop/crates/ai_bundle) for more info on the `ai app` and `ai saf` commands
- [aingle_cli_sandbox](https://github.com/AIngleLab/aingle/tree/develop/crates/ai_sandbox) for more info on the `ai sandbox` command

## Installation

### Requirements

- [Rust](https://rustup.rs/)
- [AIngle](https://github.com/AIngleLab/aingle) binary on the path

### Building

From github:

```shell
cargo install aingle_cli --git https://github.com/AIngleLab/aingle
```

From the aingle repo:

```shell
cargo install --path crates/ai
```
