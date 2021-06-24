#!/usr/bin/env bash

set -euxo pipefail
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

owner=aingle
repo=aingle
box=${1:-latest}
branch=${2:-$( git rev-parse --abbrev-ref HEAD )}
tag="$owner/$repo:$box.$branch"

docker push "aingle/aingle:$box.$branch"
