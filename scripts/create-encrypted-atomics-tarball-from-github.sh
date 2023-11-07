#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <output_path>"
    exit 1
fi

output_path="${1}"
rm -f "${output_path}"
wget -O- -q https://api.github.com/repos/redcanaryco/atomic-red-team/tarball | age -p > "${output_path}"
