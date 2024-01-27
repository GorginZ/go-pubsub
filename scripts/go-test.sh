#!/usr/bin/env bash
set -eo pipefail

usage() {
    echo "Usage: $0 <module_to_test>"
    echo "  order - test order"
    echo "  payment - test payment"
    echo "  shipping - test shipping"
    echo "  all - test all"
    exit 1
}

[[ $# -eq 1 ]] || usage

die() { echo "$1" >&2; exit "${2:-1}"; }

hash docker || die "docker not found" $?

module_to_test=$1
case $module_to_test in
    order|payment|shipping)
        docker compose run --rm go-sh sh -c "cd go-pubsub/$module_to_test && go test -v ./...";
        ;;
#example for running all tests for the exercise, we'll build these services out too
    all)
        docker compose run --rm go-sh sh -c "cd go-pubsub/order && go test -v ./... && cd ../payment && go test -v ./... && cd ../shipping && go test -v ./..."
        ;;

    *)
        echo "invalid action"
        exit 1
        ;;
esac
