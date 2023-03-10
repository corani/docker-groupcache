#!/bin/bash 

function print_usage {
    echo "Usage: $(realpath "$0") [options]"
    echo
    echo "Available options:"
    echo "  -h      prints this help"
    echo "  -b      local build"
    echo "  -d      docker build"
    echo "  -up     docker-compose up"
    echo "  -down   docker-compose down"
}

function local_build {
    do_echo go build -o app .
}

function docker_build {
    do_echo docker build -t corani/docker-groupcache .
}

function docker_up {
    do_echo docker-compose up
}

function docker_down {
    do_echo docker-compose down
}

function do_echo {
    echo "[CMD] $@"
    TIMEFORMAT="[INFO] took %3lR"
    time "$@"
    code=$?
    if [ $code -ne 0 ]; then
        rc=$code
        echo "[ERROR] return code $rc"
    fi
}

while [ "$#" -gt "0" ]; do
  arg=$1
  shift

  case $arg in
    -h)
        print_usage
        exit 0
        ;;
    -b)
        local_build
        ;;
    -d)
        docker_build
        ;;
    -up)
        docker_up
        ;;
    -down)
        docker_down
        ;;
    *)
        echo "[ERROR] unrecognized argument '$arg'"
        print_usage
        exit 1
        ;;
  esac
done

