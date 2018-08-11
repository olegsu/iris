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
    echo $version
}

updateChartVersion(){
    version=$1
    yq --arg version "$version" '.version = $version' $CF_VOLUME_PATH/iris/iris/Chart.yaml | yq -y '.sources[.sources | length] = env.CF_COMMIT_URL' --yaml-output > $CF_VOLUME_PATH/Chart.new.yaml
    mv $CF_VOLUME_PATH/Chart.new.yaml $CF_VOLUME_PATH/iris/iris/Chart.yaml
}

v=$(getChartVersion)
echo "Setting version to be $v"
cat $CF_VOLUME_PATH/iris/iris/Chart.yaml