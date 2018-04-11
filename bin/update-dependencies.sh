#!/usr/bin/env bash

set -e

self_dir=$(dirname "$0")
base_dir=$( cd "$self_dir/.." ; pwd -P )

function sayPass() {
  echo -e "\033[0;32mPASS\033[0m" 
}

function sayFail() {
  echo -e "\033[0;31mFAIL\033[0m"
}

cd "$base_dir"
go get -t -u -v
go test -v && sayPass || sayFail

