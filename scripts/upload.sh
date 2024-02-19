#!/bin/bash
set -euxo pipefail

help_message="Usage: upload.sh <ASSETS_DIR> <RELEASE_TAG>"

assets_dir=${1:-}
release_tag=${2:-}

if [ ${#assets_dir} -eq 0 -o ${#release_tag} -eq 0 ]; then
  echo "${help_message}"
  exit 1
fi

assets=($(find "${assets_dir}" -type f))
if [ "${#assets[@]}" -eq 0 ]; then
  echo "error: no files found in ${assets_dir}" >&2
  exit 1
fi

gh release upload "${release_tag}" --clobber -- "${assets[@]}"
