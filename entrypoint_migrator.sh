#!/bin/bash

DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME port=$DBPORT"

goose postgres "$DBSTRING" up;
