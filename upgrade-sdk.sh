#!/bin/bash

CMD="GOPRIVATE=github.com/myscribae/myscribae-sdk-go go get github.com/myscribae/myscribae-sdk-go"

echo "$CMD"
eval "$CMD"