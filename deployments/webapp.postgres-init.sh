#!/bin/bash

set -e
set -u

function create_user_and_database() {
	local database=$1
	local user=$2
	local password=$3
	echo "  Creating user and database '$database'"
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
      CREATE USER $user WITH PASSWORD '$password';
	    CREATE DATABASE $database;
	    GRANT ALL PRIVILEGES ON DATABASE $database TO $user;
EOSQL
}

create_user_and_database "$HYDRA_DB" "$HYDRA_USERNAME" "$HYDRA_PASSWORD"
create_user_and_database "$CORE_DB" "$CORE_USERNAME" "$CORE_PASSWORD"
