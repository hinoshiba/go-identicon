#!/bin/sh
export GOPATH
GOPATH="`pwd`"
cd src/iDicon
dep ensure
dep status
