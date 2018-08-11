#!/bin/bash

defaultBranch="master"

getChartVersion(){
    version="$(cat VERSION)"
    if [ "$CF_BRANCH_TAG_NORMALIZED" = "$defaultBranch" ]
    then
        version=$version
    else
        version=$version-$CF_BRANCH_TAG_NORMALIZED-$CF_SHORT_REVISION
    fi
    echo "Setting version to be $version"
    echo $version
}

echo $getChartVersion