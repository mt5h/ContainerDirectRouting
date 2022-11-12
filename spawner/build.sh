#!/usr/bin/env bash

set -e

export GOOS=${GOOS:-$(go env GOOS)}
export GOARCH=${GOARCH:-$(go env GOARCH)}
export CGO_ENABLED=${GO_CGO_ENABLED:-"1"}
GO_FLAGS=${GO_FLAGS:-"-tags netgo"}    # Extra go flags to use in the build.
BUILD_USER=${BUILD_USER:-"${USER}@${HOSTNAME}"}
BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}
VERBOSE=${VERBOSE:-}
OUTPUT_NAME_WITH_ARCH=${OUTPUT_NAME_WITH_ARCH:-"false"}

repo_path_spawner="ContainerDirectRouting/spawner"

version=${VERSION:-$( git describe --tags --dirty --abbrev=14 | sed -E 's/-([0-9]+)-g/.\1+/' )}
revision=$( git rev-parse --short HEAD 2> /dev/null || echo 'unknown' )
branch=$( git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown' )
go_version=$( go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/' )


# go 1.4 requires ldflags format to be "-X key value", not "-X key=value"
ldseparator="="
if [ "${go_version:0:3}" = "1.19" ]; then
	ldseparator=" "
fi

ldflags_spawner="
  -X ${repo_path_spawner}/version.Version${ldseparator}${version}
  -X ${repo_path_spawner}/version.Revision${ldseparator}${revision}
  -X ${repo_path_spawner}/version.Branch${ldseparator}${branch}
  -X ${repo_path_spawner}/version.BuildUser${ldseparator}${BUILD_USER}
  -X ${repo_path_spawner}/version.BuildDate${ldseparator}${BUILD_DATE}
  -X ${repo_path_spawner}/version.GoVersion${ldseparator}${go_version}"


if [ -n "$VERBOSE" ]; then
  echo "Building with -ldflags $ldflags_spawner"
fi

mkdir -p "$PWD/_output"
output_file="$PWD/_output/spawner"
if [ "${OUTPUT_NAME_WITH_ARCH}" = "true" ] ; then
  output_file="${output_file}-${version}-${GOOS}-${GOARCH}"
fi

pushd spawner > /dev/null
go build ${GO_FLAGS} -ldflags "${ldflags_spawner}" -o "${output_file}" "$PWD"
popd > /dev/null

exit 0
