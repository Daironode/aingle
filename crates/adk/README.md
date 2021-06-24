# AIngle Development Kit (ADK)

This kit:

1. The DSL is ergonomic and composable, so optional if you want more control
2. Differentiates between the aingle API/interface and "sugar" syntax


## ADK API

There are low-level macros and high-level functions to aid writing happs.

The intention is that most of the time most developers will use the high level
functions as they leverage the Rust type system better than macros can. This
allows for more useful compiler and IDE feedback loops.

## Examples

### adk_extern

```rust
use crate::prelude::*;

#[derive(Debug, Serialize, Deserialize, SerializedBytes)]
pub struct MyInput;

#[derive(Debug, Serialize, Deserialize, SerializedBytes)]
pub struct MyOutput(MyInput);

#[adk_extern]
fn foo(input: MyInput) -> ExternResult<MyOutput> {
  Ok(MyOutput(input))
}

```

### adk_entry

```rust

#[adk_entry(id = "foo")]
#[derive(Clone)]
pub struct Foo;

#[adk_entry(id = "bar")]
#[derive(Clone)]
pub struct Bar;

entry_defs![Foo::entry_def(), Bar::entry_def()];
```

### create_entry, get, hash_entry, create_link, get_links, debug!

```rust
// Create your entry types
let foo = Foo;
let bar = Bar;
// Commit the entries
let _foo_header_hash = create_entry(foo.clone())?;
let _bar_header_hash = create_entry(bar.clone())?;
// Get the entry hash of each entry
let foo_entry_hash = hash_entry(foo)?;
let bar_entry_hash = hash_entry(bar)?;
// Link from foo (base) to bar (target)
let _link_add_header_hash = create_link(foo_entry_hash.clone(), bar_entry_hash)?;
// Get the links back
let links = get_links(foo_entry_hash)?;
// Print out the links
debug!(links);
```

### call_remote, zome_info, agent_info

```rust
// Get your agent key
let agent_pubkey = agent_info()?.agent_pubkey;
// Get the name of this zome
let zome_name = zome_info()?.zome_name;
// Call your friends "foo" function
let result: SerializedBytes = call_remote(
    my_friends_agent_pubkey,
    zome_name,
    "foo".to_string(),
    CapSecret::default(),
    MyInput
)?;
// Get their output
let output: MyOutput = result.decode()?;
// Print their output
debug!(output);
```

### Direct Api Call
The above macros are convenience macros for calling the api but this
can also be done directly as follows:

```rust
// Commit foo
let foo_header_hash = create_entry(foo.clone())?;
// Call the api directly:
// Create the Entry from bar.
let entry = Entry::App(bar.clone().try_into()?);
// Call the update_entry host_fn directly
let _bar_header_hash = host_call::<UpdateInput, HeaderHash>(
    __update,
    UpdateInput::new(foo_header_hash, EntryWithDefId::new(Bar::entry_def().id, entry)))
)?;
```

Current version: 0.0.100

## Composable concepts

One of the main design goals of this ADK is to make it composable.

The macros mostly just remove boilerplate that tends to obfuscate core aingle
concepts and allow for bugs to sneak in.

These macros are designed to be largely "mechanical" though, if you want to do
something a little bespoke then there is always a more verbose option to fall
back on.

This means the abstractions provided by the ADK are optional and composable so
a developer can opt in to only the functionality that is useful to their app.

The test wasms used by aingle core are written without the ADK.

The _mandatory_ wasm components are instead maintained in the `aingle-wasmer`
crate https://github.com/AIngleLab/aingle-wasmer. These mandatory components
exist because there needs to be a basic protocol that aingle can implement
to co-ordinate memory and callbacks with the wasm.

The `aingle-wasmer` repository contains 3 main crates:

- `aingle_wasmer_common`: shared abstractions for both the host and guest
- `aingle_wasmer_host`: implements wasm for aingle itself
- `aingle_wasmer_guest`: abstractions for _you_ to write wasm with, that also
  power the ADK under the hood

It is important that it is possible to write minimal wasms that are compatible
with aingle without pulling in "the kitchen sink" of irrelevant Rust
dependencies or hiding so many details behind a DSL that developers really have
no idea what is going on and end up "cargo culting" solutions wholesale.

The `aingle-wasmer` crate has its own detailed documentation but relevant
high level details will be included here.

## AIngle overview

AIngle has several high-level components:

- A SGD network that shares, validates and stores data
- Wasm & SAF files that are executed to provide application specific logic
- A user-facing websockets interface that enables interactive clients
- The aingle binary that co-ordinates all these components

Depending on which component(s) you are working with, the key concepts and
documentation may look very different.

This documentation describes how to write wasm files that are compatible with
the aingle core binary.

If this is your first time writing wasm, or even Rust code, don't worry!

The ocean of wasm and Rust development is vast and deep, but you only need to
dip your toes in to effectively write wasm for aingle.

- AIngle core handles many of the tough edge-cases for you, like checking
  cryptographic proofs and detecting common "bad behaviour" on the network
- The ADK (aingle development kit) provides a DSL (domain specific language)
  to remove most or all boilerplate
- Most of the low-level wasm limitations have been abstracted away, so you can
  mostly just write vanilla rust, using all the standard language features
- Most of the advanced functionality in Rust is not required, there is
  little or no need for multithreading, channels, locks, complex traits,
  lifetimes, etc. etc.

Every aingle wasm works in the same basic way. The application developer
writes some Rust code using the functionality exposed by aingle. As long as
the rust code can be compiled to wasm and exposes the interface that aingle
expects, then aingle can run it to manage a p2p SGD network.

**There are three things that make a wasm aingle-compatible:**

- **It must use only the host functionality that aingle provides**
- **It must expose callback functions that aingle expects**
- **Memory handling and (de)serialization must be compatible with aingle**

### AIngle functionality

AIngle exposes a list of aingle-specific things that a wasm can do.

For detailed documentation of the full list, see the `core/ribosome` module
inside core, but some illustrative examples include:

- `emit_signal`: publish data to subscribed interface clients
- `sign`: use the agent's keypair to sign some data
- `create_entry`: save some data to the local source chain and broadcast it to
  the SGD to be redundantly validated and stored
- `get_entry`: retrieve some data from local or the network given its hash
- `create_link`: create graph style relationships (links) between entries
- `get_links`: retrive links between entries using the SGD as a graph database
- `remote_signal`: send data directly to known peers on the network without waiting for a response
- `call_remote`: perform a remote procedure call on a peers node, if you're authorized

This toolkit of functionality is available to the wasm as a list of "extern"
functions that are all injected into the wasm by aingle - i.e. these
functions are all provided by aingle to be used by every wasm.

All of this functionality is enabled on the wasm guest by the
`aingle_externs!` macro in the `aingle_wasmer_guest` crate. It just needs
to be called once somewhere in the wasm.

### AIngle interface

AIngle drives all of the wasms it has installed in the same way.

Internally aingle handles all the multi-threading, co-ordination between the
network and websocket RPC connections to any interactive client (e.g. like an
 electron app).

Whenever aingle reaches some point where it needs to execute application
specific logic it will call one of the functions in the wasm directly.

This is analagous to how standalone binaries in Rust (and other languages)
start by running the `main` function by convention.

It is also similar to how "handler" functions are configured in common
"serverless" platforms like AWS lambda.

There are broadly two types of exposed functions:

- well defined callbacks that specific functionality e.g. "validate this entry"
- arbitrary functions that handle serialized data from the interactive client

Both work in the same way technically but they have different responsibilities.
Callbacks extend aingle itself whereas other extern functions extend the
interactive client that is sending requests via aingle.

All the extern functions are run in a newly built, sandboxed wasm instance, so
there are no long-running processes and it is not possible for callbacks to
interact with each other directly or share data in memory.

Because all functionality is based on simple, sandboxed callbacks, there is no
need for the application developer to handle threading or other complexities.
Even minor memory leaks are relatively harmless as the wasm memory is dropped
wholesale after every extern function call.

### AIngle serialization and memory

Due to wasm limitations (see below) aingle must implement some process to
share complex data types between the host and the guest.

The full process is documented in detail in the `aingle-wasmer` repository.

https://github.com/AIngleLab/aingle-wasmer

In short, there are a few functions that the guest needs to expose to the host
that the host will use to request safe memory allocations and deallocations from
the guest.

This allows the host to repect the guest's own memory allocation logic, and so
provides support for alternative allocators.

Exposing these functions is as simple as calling the `aingle_externs!` macro
in the `aingle_wasmer_guest` crate.

Once the host and guest can share memory safely, they need to decide on a
serialization format that data can be shared across the wasm boundary as.

AIngle uses the messagepack serialization format as it has several benefits:

- It is not tied to the rust compiler
- It is reasonably fast and compact
- It supports binary data natively (e.g. JSON does not)
- It is reasonably human readable and can even be automatically JSONified

If you aren't familiar with `serde`, messagepack and/or
`aingle_middleware_bytes`. then it's worth at least skim reading the
documentation.

- Messagepack: https://msgpack.org/index.html
- Serde: https://github.com/serde-rs/serde
- Serde messagepack: https://github.com/3Hren/msgpack-rust
- AIngle serialized bytes: https://github.com/AIngleLab/aingle-serialization/tree/develop/crates/aingle_middleware_bytes

Check out the [advanced AIngle + WASM Overview for deeper details](./ON-WASM.md).
