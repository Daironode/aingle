# ai_bundle

[![Project](https://img.shields.io/badge/project-aingle-blue.svg?style=flat-square)](http://aingle.ai/)
[![Forum](https://img.shields.io/badge/chat-forum%2eaingle%2enet-blue.svg?style=flat-square)](https://forum.aingle.ai)
[![Chat](https://img.shields.io/badge/chat-chat%2eaingle%2enet-blue.svg?style=flat-square)](https://chat.aingle.ai)

[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

Utilities to create SAF and hApp bundle files from a source working directory, and vice-versa.
This crate defines two separate subcommands for the `ai` CLI tool, one for each type of bundle.
Both subcommands are very similar and have identical interfaces.

This crate also defines standalone binaries for each subcommand, `ai-saf` and `ai-app`.

Usage instructions from the `-h` flag:

```sh
$ ai saf -h

ai-saf 0.0.1
Work with SAF bundles

USAGE:
    ai saf <SUBCOMMAND>

FLAGS:
    -h, --help       Prints help information
    -V, --version    Prints version information

SUBCOMMANDS:
    help      Prints this message or the help of the given subcommand(s)
    init      Create a new, empty AIngle SAF bundle working directory
    pack      Pack the contents of a directory into a `.saf` bundle file
    unpack    Unpack the parts of `.saf` file out into a directory
```

`ai app -h` is very similar.

## Contribute


* Connect with us on our [forum](https://forum.aingle.ai)

## License
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

Copyright (C) 2019 - 2021, AIngle

This program is free software: you can redistribute it and/or modify it under the terms of the license
provided in the LICENSE file (Apache 2.0).  This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
PURPOSE.
