#!/bin/bash

# Fail on any error
set -e


# install oapi-codegen 
if ! command -v oapi-codegen &> /dev/null
then
    echo "oapi-codegen could not be found, installing..."
    go install "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest"
else
    echo "oapi-codegen is already installed."
fi

echo "Generating Go code from OpenAPI spec..."

# Run oapi-codegen for types and server generation
oapi-codegen -generate chi-server,types,spec -package api -o api/api.gen.go api/openapi.yaml

echo "Code generation completed successfully!"