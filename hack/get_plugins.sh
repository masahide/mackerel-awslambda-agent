#!/bin/bash

set -x
WORKDIR="$TMPDIR$(date +%s)"

URL="https://github.com/mackerelio-labs/check-aws-cloudwatch-logs-insights/releases/download/v0.0.2/check-aws-cloudwatch-logs-insights_linux_amd64.zip"
FILE="check-aws-cloudwatch-logs-insights"

mkdir $WORKDIR

curl -Ls \
    -o "$WORKDIR/plugin.zip" \
    "$URL"
unzip "$WORKDIR/plugin.zip" -d "$WORKDIR" 
mv "$WORKDIR/$FILE"*/"$FILE" .
rm -rf "$WORKDIR"

