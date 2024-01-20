#!/usr/bin/env bash
set -eo pipefail

usage() {
    echo "Usage: $0 <ci_action>"
    echo "  go - lint go code"
    echo "  tf - lint terraform code"
    exit 1
}

[[ $# -eq 1 ]] || usage

die() { echo "$1" >&2; exit "${2:-1}"; }

hash docker || die "docker not found" $?

ci_action=$1

lint_go() {
    echo "linting go"
    docker compose run --rm go-lint run -v go-pubsub/order/*.go
}

#fmt and validate tf
tf_fmt_and_validate() {
    tf_single_action tf_fmt
    tf_single_action tf_validate
}

tf_single_action() {
    case $ci_action in
        tf_fmt)
            docker compose run --rm terraform fmt -check -recursive || die "terraform fmt failed" $?
            ;;
        tf_validate)
            docker compose run --rm terraform validate -no-color || die "terraform validate failed" $?
            ;;
        *)
            echo "invalid action"
            exit 1
            ;;
    esac
}

case $ci_action in
    lint_go)
        lint_go
        ;;
    tf_fmt|tf_validate)
        tf_single_action
        ;;
    validate_all)
        lint_go
        tf_fmt_and_validate
        ;;
    *)
        echo "invalid action"
        exit 1
        ;;
esac
