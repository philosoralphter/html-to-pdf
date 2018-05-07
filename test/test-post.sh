#!/usr/bin/env bash

curl -d "@test-html-body.html" -X POST http://localhost:9190/convert -o ../deployment/output/test-output.pdf

