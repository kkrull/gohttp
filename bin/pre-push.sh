#!/usr/bin/env bash

set -e

self_dir=$(dirname "$0")
base_dir=$( cd "$self_dir/.." ; pwd -P )

cd "$base_dir"
git diff-index --quiet HEAD --
if (( $? != 0 ))
then
  echo "Uncommitted changes"
  exit 1
fi

goimports -w .
git diff-index --quiet HEAD -- || git commit -am "Formatting"
