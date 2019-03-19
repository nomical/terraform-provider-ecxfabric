#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

DIST_PATH=$1
VERSION=$2

gox -os="linux darwin windows" -arch="amd64" -output="$DIST_PATH/{{.OS}}_{{.Arch}}"

cd $DIST_PATH

for f in $(find . -maxdepth 1 -type f)
do
    b="$(basename $f)"
    d="d_$b"
    mkdir "$d"

    if [[ $b == windows* ]]
    then
        mv $b "$d/terraform-provider-ecxfabric_${VERSION}.exe"
        cd $d
        zipName=${b%".exe"}
        find . -maxdepth 1 -type f -exec zip -m "$zipName.zip" "{}" \;
    else
        mv $b "$d/terraform-provider-ecxfabric_${VERSION}"
        cd $d
        find . -maxdepth 1 -type f -exec tar --remove-files -czf "$b.tar.gz" {} \;
    fi
    
    mv * ../.
    cd ..
done

find . -maxdepth 1 -mindepth 1 -type d -exec rm -rf "{}" \;
