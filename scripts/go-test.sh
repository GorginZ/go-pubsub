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
        docker compose run --rm go-sh sh -c "cd go-pubsub/$test_action && go test ./...";
        ;;
#example for running all tests for the exercise, we'll build these services out too
    all)
        docker compose run --rm go-sh sh -c "cd go-pubsub/order && go test ./... && cd ../payment && go test ./... && cd ../shipping && go test ./..."
        ;;

    *)
        echo "invalid action"
        exit 1
        ;;
esac
