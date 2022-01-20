#!/bin/bash

# generate swagger spec file

# precondition： 
# - install swagger cli first : https://goswagger.io/install.html

set -e

function genPublicSpec() {
  echo "Generating public spec..."
  swagger generate spec \
    --include-tag=public \
    -i ./api/swagger_ui/api-base-public.yml \
    -o api/swagger_ui/api.yml
  echo "Generated public spec"
}

function genInnerSpec() {
  echo "Generating inner spec..."
  swagger generate spec \
  -i ./api/swagger_ui/api-base-inner.yml \
  -o api/swagger_ui/api-inner.yml

  swagger generate markdown -f ./api/swagger_ui/api.yml --output=api.md

  echo "Generated inner spec"
}


case $1 in
  '')
    genPublicSpec
    genInnerSpec
    ;; 
	public )
		genPublicSpec ;;
	inner )
		genInnerSpec ;;
	* )
		echo "unknown subcommand"
		exit 1 
		;;
esac