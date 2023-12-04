#!/usr/bin/env bash

make

make test

INPUT_FILES=$(find ./res/ -type f -name "test*input.csv")

for file in $INPUT_FILES
do
    echo "Testing $file"
    ./bin/program input $file
done
