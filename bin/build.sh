#!/usr/bin/env bash
set -e

dingo -src=./internal/config/di -dest=./generated

docker build . -t navexplorer/indexer:dev
docker push navexplorer/indexer:dev
