#!/bin/bash

DEFAULT_DATA_DIR='./.data'
DEFAULT_TAG='latest'

help() {
    echo "Runs the HITS hit counter API and database"
    echo "Usage: run.sh [DATABASE_PASSWORD] [SALT] [HONEYCOMB_KEY] [DATA_DIRECTORY]"
    echo "  -p DATABASE_PASSWORD - Password to be used to communicate to the database. Required."
    echo "  -s SALT              - Salt used to anonymize user data. Required."
    echo "  -k HONEYCOMB_KEY     - Key to send data to Honeycomb. Optional."
    echo "  -d DATA_DIRECTORY    - Local directory that the database will use for persistence. Default: ${DEFAULT_DATA_DIR}. Optional."
    echo "  -i INTERACTIVE       - Runs the services in interactive mode. Optional."
    echo "  -t TAG               - Image tag to target. Default: ${DEFAULT_TAG}. Optional."
    echo "  -r REGISTRY          - Registry to use to pull the image"
}

if [ $# -eq 0 ]
  then
    echo "No arguments provided"
    help
    exit 1
fi

while getopts p:s:k:d:r:hit: OPTION;
do
    case $OPTION in
    p)
        DATABASE_PASSWORD=$OPTARG
        ;;
    s)
        SALT=$OPTARG
        ;;
    k)
        HONEYCOMB_KEY=$OPTARG
        ;;
    d)
        DATA_DIRECTORY=$OPTARG
        ;;
    i)
        INTERACTIVE=true
        ;;
    t)
        TAG=$OPTARG
        ;;
    r)
        REGISTRY=$OPTARG
        echo $OPTARG
        ;;
    h)
        help
        exit 0
        ;;
    \?)
        exit 1
        ;;
    :)
        exit 1
        ;;
    *)
        help
        exit 1
        ;;
    esac
done

if [ -z "${DATABASE_PASSWORD}" ] || [ -z "${SALT}" ]; then
    help
    exit 1
fi

# Assume local if no parameter is provided
if [[ -z ${DATA_DIRECTORY} ]]
then 
    echo "No data directory provided. Using ${DEFAULT_DATA_DIR}"
    DATA_DIR=${DEFAULT_DATA_DIR}
else
    DATA_DIR=${DATA_DIRECTORY}
fi

echo "Using '${DATA_DIR}' as the database data directory"
echo "DATA_DIR=${DATA_DIR}" > .env
echo "DATABASE_PASSWORD=${DATABASE_PASSWORD}" >> .env
echo "SALT=${SALT}" >> .env
echo "HC_KEY=${HONEYCOMB_KEY}" >> .env

if [[ -z ${TAG} ]]
then 
    echo "No tag provided. Using ${DEFAULT_TAG}"
    echo "TAG=${DEFAULT_TAG}" >> .env
else
    echo "TAG=${TAG}" >> .env
fi

echo "Using registry: ${REGISTRY}"
echo "REGISTRY=${REGISTRY}" >> .env

if [[ -z ${INTERACTIVE} ]]
then
    docker-compose up -d
else
    docker-compose up
fi