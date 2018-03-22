#!/usr/bin/env bash

set -e

self_dir=$(dirname "$0")
base_dir=$( cd "$self_dir/.." ; pwd -P )

cd "$base_dir"
go build
./gohttp $@ \
  && echo -e "\033[0;32mPASS\033[0m" \
  || echo -e "\033[0;31mFAIL\033[0m"
