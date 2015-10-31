#!/usr/bin/env bats

load helpers

function setup() {
    _setup
    gattai init
    cat > provision.yml << EOF
---
machines:
EOF
}

function teardown() {
    _teardown
}

@test "gattai add -f not-existed" {
    run gattai add --flavor not-existed test
    [ "$status" -ne 0 ]
}

@test "gattai add --flavor none" {
    run gattai add --flavor none test
    [ "$status" -eq 0 ]

    run gattai provision test
    echo $output
    [ "$status" -eq 0 ]

    run gattai active test
    [ "$status" -eq 0 ]

    run stat .gattai/.active_host
    [ "$status" -eq 0 ]

    run cat .gattai/.active_host
    [ "${lines[1]}" = "name: test" ]
    [ "${lines[2]}" = "DOCKER_HOST: \"tcp://1.2.3.4:2376\"" ]
}