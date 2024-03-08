#!/bin/bash

./startup.sh up -s couchdb && ./startup.sh createChannel && ./ccInstall full
