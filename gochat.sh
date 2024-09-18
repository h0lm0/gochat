#!/bin/bash

startGochat(){
    if [ -n "$ENV" ]; then
        docker compose -f docker-compose."${ENV}".yml --env-file ./app/.$ENV.env up -d --build
    else
        echo "Please use -e to define environment ('dev' | 'prod') before -u"
    fi
}

stopGochat(){
    if [ -n "$ENV" ]; then
        docker compose -f docker-compose."${ENV}".yml --env-file ./app/.$ENV.env down
    else
        echo "Please use -e to define environment ('dev' | 'prod') before -d"
    fi
}

OPTSTRING=":ude:"

while getopts ${OPTSTRING} opt; do
  case ${opt} in
    u)
      startGochat
      ;;
    d)
      stopGochat
      ;;
    e)
      ENV=$OPTARG
      ;;
    ?)
      echo "Invalid option: -${OPTARG}."
      exit 1
      ;;
    :)
      echo "Option -${OPTARG} requires an argument."
      exit 1
      ;;
  esac
done
