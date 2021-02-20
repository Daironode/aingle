{

 # configure ainglenix itself
 ainglenix = {

  # true = use a github repository as the ainglenix base (recommended)
  # false = use a local copy of ainglenix (useful for debugging)
  use-github = true;

  # configure the remote ainglenix github when use-github = true
  github = {

   # can be any github ref
   # branch, tag, commit, etc.
   ref = "6decb861920a517fae5840acf79766196ede9b0d";

   # the sha of what is downloaded from the above ref
   # note: even if you change the above ref it will not be redownloaded until
   #       the sha here changes (the sha is the cache key for downloads)
   # note: to get a new sha, get nix to try and download a bad sha
   #       it will complain and tell you the right sha
   sha256 = "0f6p2vz32g8fgphkvfiswlg3wb0bpld6nsagfpmc2ad1g0l85xqy";

   # the github owner of the ainglenix repo
   owner = "aingle";

   # the name of the ainglenix repo
   repo = "ainglenix";
  };

  # configuration for when use-github = false
  local = {
   # the path to the local ainglenix copy
   path = ../ainglenix;
  };

  includeAIngleBinaries = false;
 };

 release = {
  hook = {
   # sanity checks before deploying
   # to stop the release
   # exit 1
   preflight = ''
hn-release-hook-preflight-manual
'';

   # bump versions in the repo
   version = ''
hn-release-hook-version-rust
hcp-release-hook-version
'';

   # publish artifacts to the world
   publish = ''
# crates are published from circle!
'';
  };

  # the commit hash that the release process should target
  # this will always be behind what ends up being deployed
  # the release process needs to add some commits for changelog etc.
  commit = "8fb82a3a6d8cc69c95c654bd21bf15785a6ca291";

  # the semver for prev and current releases
  # the previous version will be scanned/bumped by release scripts
  # the current version is what the release scripts bump *to*
  version = {
   current = "0.0.3";
   # not used by version hooks in this repo
   previous = "_._._";
  };

  github = {
   # markdown to inject into github releases
   # there is some basic string substitution {{ xxx }}
   # - {{ changelog }} will inject the changelog as at the target commit
   template = ''
{{ changelog }}

# Installation

Use AInglenix to work with this repository.

See:

- https://github.com/Daironode/ainglenix
- https://nixos.org/
'';

   # owner of the github repository that release are deployed to
   owner = "aingle";

   # repository name on github that release are deployed to
   repo = "aingle";

   # canonical local upstream name as per `git remote -v`
   upstream = "origin";
  };
 };
}