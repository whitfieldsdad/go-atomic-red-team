#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <input_dir> <output_path>"
    exit 1
fi

input_dir=$1
output_path=$2

tar czvf "${output_path}" "${input_dir}"
