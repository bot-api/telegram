#!/bin/bash

set -o errexit
set -o pipefail

main() {
  _cd_into_top_level
  _generate_coverage_files
}

_cd_into_top_level() {
  cd "$(git rev-parse --show-toplevel)"
}

_generate_coverage_files() {
  for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d); do
    if ls $dir/*.go &>/dev/null ; then
      go test -tags integration -run TestI_* -covermode=count -coverprofile=$dir/profile_i.coverprofile $dir || fail=1
    fi
  done
}

main "$@"