#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $0); pwd)
URL="https://github.com/mackerelio-labs/check-aws-cloudwatch-logs-insights/releases/download/v0.0.2/check-aws-cloudwatch-logs-insights_linux_amd64.zip"
FILE="check-aws-cloudwatch-logs-insights"


[[ ${DIST} == "" ]] && DIST="${SCRIPT_DIR}/../.dist/checker/"


get_check-aws-cloudwatch-logs-insights () {

    WORKDIR="$TMPDIR$(date +%s)"
    mkdir $WORKDIR
    curl -Ls \
    -o "$WORKDIR/plugin.zip" \
    "$URL"  \
    && unzip "$WORKDIR/plugin.zip" -d "$WORKDIR" 
    [[ -d $DIST ]] || mkdir -p $DIST
    mv "$WORKDIR/$FILE"*/"$FILE" $DIST/
    rm -rf "$WORKDIR"
}

set -x
[[ -f ${DIST}${FILE} ]] || get_check-aws-cloudwatch-logs-insights
