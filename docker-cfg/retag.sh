#!/bin/bash

set -e

CONFIG_FILE=$1
SOURCE_TAG=$2
NEW_TAG=$3
WORKING_DIR=$(dirname $0)

python $WORKING_DIR/yaml_tags.py $1 | xargs -I % docker tag %:$SOURCE_TAG %:$NEW_TAG
