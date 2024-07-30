#!/bin/bash

DBSTRING="host=$PG_HOST user=$PG_USER password=$PG_PASSWORD dbname=$PG_DATABASE port=$PG_PORT"

goose postgres "$DBSTRING" up;
