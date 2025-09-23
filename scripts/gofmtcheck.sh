#!/bin/bash

# Check goimports
echo "==> Checking that code complies with goimports requirements..."
goimports_files=$(goimports -l ./edgenext)
if [[ -n ${goimports_files} ]]; then
    echo 'goimports needs running on the following files:'
    echo "${goimports_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."
gofmt_files=$(gofmt -l ./edgenext)
if [[ -n ${gofmt_files} ]]; then
    echo 'gofmt needs running on the following files:'
    echo "${gofmt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

exit 0
