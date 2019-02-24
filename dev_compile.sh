#!/bin/sh
NOW=`date '+%Y%m%d%H%M%S'`
export GOPATH="`pwd`"
ls -1 src/iDicon/exec | while read row ; do
  GOOS=linux GOARCH=amd64 go install -ldflags "-s -w" iDicon/exec/$row
  GOOS=windows GOARCH=amd64 go install -ldflags "-s -w" iDicon/exec/$row
  GOOS=darwin GOARCH=amd64 go install -ldflags "-s -w" iDicon/exec/$row
done
