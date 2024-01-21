#!/usr/bin/env bash
set -eo pipefail

usage() {
    echo "Usage: $0 <test_action>"
    echo "  order - test order"
    echo "  payment - test payment"
    echo "  shipping - test shipping"
    echo "  all - test all"
    exit 1
}

[[ $# -eq 1 ]] || usage

die() { echo "$1" >&2; exit "${2:-1}"; }

hash docker || die "docker not found" $?

test_action=$1

case $test_action in
    order|payment|shipping)
        # run test for specific service only no need to find
        # docker compose run --rm go-sh cd go-pubsub/$test_action && go test ./...  ## todo fix this difficulties with the compose service entrypoint for the alpine go img.
        ;;

    all)
        # docker compose run --rm go-sh find . -name go.mod -execdir go test ./... \; #this cmd works locally but alpine uses busybox find and doesnt have execdir...
        docker compose run --rm go-sh find . -name go.mod -exec sh -c 'cd $(dirname {}); go test ./...' \; | grep -i fail && die "tests failed" 1 # we have to exit here because the way we do all tests it exits 0 and so does grep on match. bad, I know
        ;;

    *)
        echo "invalid action"
        exit 1
        ;;
esac
