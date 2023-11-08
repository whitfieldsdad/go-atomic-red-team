#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <output_path>"
    exit 1
fi

output_path=$1
rm -f "${output_path}"

url="https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json"
wget -O- -q "${url}" > "${output_path}"
