#! /usr/bin/env bash

set -e -x -u

pushd irc-notification-resource

  go get github.com/onsi/ginkgo/ginkgo

  make

popd
