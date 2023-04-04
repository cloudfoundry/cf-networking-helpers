#!/bin/bash

set -eu
set -o pipefail

cd $(dirname $0)/..

BIN_DIR="${PWD}/bin"
mkdir -p "${BIN_DIR}"
export PATH="${PATH}:${BIN_DIR}"

function bootDB {
  db=$1

  if [ "$db" = "postgres" ]; then
    launchDB="(docker-entrypoint.sh postgres &> /var/log/postgres-boot.log) &"
    testConnection="psql -h localhost -U postgres -c '\conninfo' &>/dev/null"
  elif [[ "$db" == "mysql"* ]]; then
    launchDB="(MYSQL_ROOT_PASSWORD=password /entrypoint.sh mysqld &> /var/log/mysql-boot.log) &"
    testConnection="echo '\s;' | mysql -h 127.0.0.1 -u root --password='password' &>/dev/null"
  else
    echo "skipping database"
    return 0
  fi

  echo -n "booting $db"
  eval "${launchDB}"
  for _ in $(seq 1 60); do
    set +e
    eval "${testConnection}"
    exitcode=$?
    set -e
    if [ ${exitcode} -eq 0 ]; then
      echo "connection established to $db"
      return 0
    fi
    echo -n "."
    sleep 1
  done
  echo "unable to connect to $db"
  exit 1
}

go build -o "$BIN_DIR/ginkgo" github.com/onsi/ginkgo/v2/ginkgo

if [ "${1:-""}" = "" ]; then
  extraArgs=""
else
  extraArgs="${@}"
fi

if [ "${DB:-"none"}" != "none" ]; then
  if [ "${DB_PORT:-"none"}" = "none" ]; then
    bootDB "${DB}"
  fi
  ginkgo -r --race -randomize-all ${extraArgs} db/timeouts
else
  echo "skipping database"
  extraArgs="--skip-package=db ${extraArgs}"
fi

ginkgo -r -p --race -randomize-all -randomize-suites -skip-package=timeouts ${extraArgs}
