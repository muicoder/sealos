#!/bin/bash

set -e

case $(uname -m) in
amd64 | x86_64)
  ARCH=amd64
  ;;
arm64 | aarch64)
  ARCH=arm64
  ;;
*)
  echo "Unsupported architecture $ARCH"
  exit
  ;;
esac

export v030=v0.30.2 v029=v0.29.6 v028=v0.28.11 v027=v0.27.15 v026=v0.26.15 v025=v0.25.16 v024=v0.24.17
k8s_api_v02=$(
  env | grep ^v02 | sort -r | while read -r v02; do
    v=${v02#*=}
    printf 's~(k8s.io/.+)%s[.0-9]+~\\1%s~g;' "${v%.*}" "$v"
  done
)

rm -rf bin
readonly TAG=${TAG:-4.3}

code_reset() {
  git reset --hard
  sed -E "s~^VERSION.+~VERSION=$TAG~;s~^BUILD_DATE.+~BUILD_DATE=2006-01-02T15:04:05-0700~" <scripts/make-rules/common.mk >f.sed && mv f.sed scripts/make-rules/common.mk
  find . -type f -name "*.mod" | while read -r mod; do sed -E "$k8s_api_v02" <"$mod" >"$mod.sed" && mv "$mod.sed" "$mod"; done
  go get github.com/google/gnostic
  go get github.com/google/gnostic-models
  go get k8s.io/kube-openapi
  go work sync
}

code_reset
{
  go work edit -replace "k8s.io/cri-api=k8s.io/cri-api@$v025"
  pushd "staging/src/github.com/labring/image-cri-shim"
  sed -E "s~(k8s.io/.+)v0.2[.0-9]+~\1$v030~g;" <go.mod >f.sed && mv f.sed go.mod
  go mod edit -replace "k8s.io/cri-api=k8s.io/cri-api@$v025" go.mod
  go mod edit -replace "github.com/containers/image/v5=github.com/containers/image/v5@v5.30.2" go.mod
  go mod edit -replace "google.golang.org/grpc=google.golang.org/grpc@v1.64.1" go.mod
  go get -u all
  git diff go.mod | grep "^[+-]"
  popd
  BINS="image-cri-shim lvscare" make build.multiarch || true
  find bin -type f -exec ls -l {} +
}
