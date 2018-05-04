#!/bin/bash

DIR="`dirname \"$0\"`"              # relative
DIR="`( cd \"$DIR\" && pwd )`"  # absolutized and normalized

APP_SOURCE_DIR="$(dirname $(dirname "$DIR"))"

#/bin/rm -rf $DIR/output
mkdir -p $DIR/output

echo "running in $DIR"

#get go dep to manage dependencies
if [ ! -e $DIR/output/dep ]
then
    echo "Downloading Godep..."
    curl -fsSLz --progress-bar https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 -o $DIR/output/dep
    chmod +x $DIR/output/dep
else
    echo "Godep found, skipping download."
fi

#compile in a docker image
echo "Compiling Go Binary..."
docker run --rm -v "$APP_SOURCE_DIR:/go/src" -v "$DIR/output:/usr/local/bin" -w "/go/src/html-to-pdf" -e GOPATH="/go"  -e CGO_ENABLED=0 -e GOOS=linux golang:latest /bin/bash -c "/usr/local/bin/dep ensure && go build -o /usr/local/bin/html-to-pdf -a -installsuffix cgo ."


#put it all in our empty docker
echo "Building html-to-pdf Docker"
docker build -t html-to-pdf:latest .

