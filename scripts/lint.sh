#!/usr/bin/env bash
set -eo pipefail

die() { echo "$1" >&2; exit "${2:-1}"; }

hash docker || die "docker not found" $?

#fmt and validate tf
docker compose run --rm terraform fmt -check -recursive || die "terraform fmt failed" $?
docker compose run --rm terraform validate -no-color || die "terraform validate failed" $?
 
#lint go code
docker compose run --rm go-lint run -v go-pubsub/order/*.go

