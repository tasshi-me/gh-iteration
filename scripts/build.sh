#!/bin/bash
set -euxo pipefail

platforms=(
#  android-amd64 # requires CGO
  android-arm64
  darwin-amd64
  darwin-arm64
  freebsd-386
  freebsd-amd64
  freebsd-arm64
  linux-386
  linux-amd64
  linux-arm
  linux-arm64
  windows-386
  windows-amd64
  windows-arm64
)

entrypoint="./cmd/gh-iteration/"

if [ ! -e "${entrypoint}" ]; then
  echo "error: entrypoint not found ${entrypoint}" >&2
fi

IFS=$'\n' read -d '' -r -a supported_platforms < <(go tool dist list) || true

for p in "${platforms[@]}"; do
  goos="${p%-*}"
  goarch="${p#*-}"
  if [[ "${supported_platforms[*]}" != *"${goos}/${goarch}"* ]]; then
    echo "error: unsupported platform ${p}" >&2
    exit 1
  fi
  ext=""
  if [ "${goos}" = "windows" ]; then
    ext=".exe"
  fi
  GOOS="${goos}" GOARCH="${goarch}" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "dist/${p}${ext}" "${entrypoint}"
done
