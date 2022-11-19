#!/usr/bin/env bash

set -e

export GOOS=${GOOS:-$(go env GOOS)}
export GOARCH=${GOARCH:-$(go env GOARCH)}
export CGO_ENABLED=${GO_CGO_ENABLED:-"1"}
GO_FLAGS=${GO_FLAGS:-"-tags netgo"}    # Extra go flags to use in the build.
BUILD_USER=${BUILD_USER:-"${USER}@${HOSTNAME}"}
BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}

version=${VERSION:-$( git describe --tags --abbrev=0 --dirty || echo 'unknown')}
revision=$( git rev-parse --short HEAD 2> /dev/null || echo 'unknown' )
branch=$( git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown' )
go_version=$( go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/' )

ldseparator="="
if [ "${go_version:0:3}" = "1.19" ]; then
	ldseparator=" "
fi

ldflags="
  -X mock-home/build.Version${ldseparator}${version}
  -X mock-home/build.Revision${ldseparator}${revision}
  -X mock-home/build.Branch${ldseparator}${branch}
  -X mock-home/build.BuildUser${ldseparator}${BUILD_USER}
  -X mock-home/build.BuildDate${ldseparator}${BUILD_DATE}
  -X mock-home/build.GoVersion${ldseparator}${go_version}
  "

echo "Setting $ldflags"

mkdir -p "$PWD/_output"
output_file="$PWD/_output/mock-home"

go build ${GO_FLAGS} -ldflags "${ldflags}" -o "${output_file}" "$PWD"

exit 0
