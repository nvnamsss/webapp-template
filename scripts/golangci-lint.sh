#!/bin/bash
GOPATH=${PWD}
EXIT_CODE=0
cd src
for var in "$@"
do
    echo "Linting src/$var"
    PKG_LIST=$(go list $var/... | grep -v /vendor/ | grep -v migrations) ; \
    golangci-lint run ${PKG_LIST} --deadline=30m
    EXIT_CODE=$(( $EXIT_CODE + $? ))
done
cd ..
exit $EXIT_CODE
