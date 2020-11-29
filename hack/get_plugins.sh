#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $0); pwd)

${DIST:="$SCRIPT_DIR/../.dist/checks"}
WORKDIR="$TMPDIR$(date +%s)"

set -x
URL="https://github.com/mackerelio-labs/check-aws-cloudwatch-logs-insights/releases/download/v0.0.2/check-aws-cloudwatch-logs-insights_linux_amd64.zip"
FILE="check-aws-cloudwatch-logs-insights"

mkdir $WORKDIR

curl -Ls \
    -o "$WORKDIR/plugin.zip" \
    "$URL"
unzip "$WORKDIR/plugin.zip" -d "$WORKDIR" 

[[ -d $DIST ]] || mkdir -p $DIST
mv "$WORKDIR/$FILE"*/"$FILE" $DIST/
rm -rf "$WORKDIR"

