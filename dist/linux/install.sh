#!/bin/bash

COMPONENT_NAME=tenantctl
SERVICE_ENV=tac.env

# Check OS and VERSION
OS=$(cat /etc/os-release | grep ^ID= | cut -d'=' -f2)
temp="${OS%\"}"
temp="${temp#\"}"
OS="$temp"

# read .env file
echo PWD IS $(pwd)
if [ -f ~/$SERVICE_ENV ]; then
    echo Reading Installation options from $(realpath ~/$SERVICE_ENV)
    env_file=~/$SERVICE_ENV
elif [ -f ../$SERVICE_ENV ]; then
    echo Reading Installation options from $(realpath ../$SERVICE_ENV)
    env_file=../$SERVICE_ENV
fi
if [ -z $env_file ]; then
    echo "No .env file found"
    exit 1
else
    source $env_file
    env_file_exports=$(cat $env_file | grep -E '^[A-Z0-9_]+\s*=' | cut -d = -f 1)
    if [ -n "$env_file_exports" ]; then eval export $env_file_exports; fi
fi

echo "Installing Tenant CLI..."

BIN_PATH=~/.local/bin
CONFIG_PATH=~/.config/$COMPONENT_NAME
CONFIG_FILE_PATH=$CONFIG_PATH/config.yaml
LOG_PATH=$CONFIG_PATH/logs

for directory in $BIN_PATH $CONFIG_PATH $LOG_PATH; do
  # mkdir -p will return 0 if directory exists or is a symlink to an existing directory or directory and parents can be created
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo "Cannot create directory: $directory"
    exit 1
  fi
done

chmod 700 $BIN_PATH
chmod 700 $CONFIG_PATH
chmod 700 $LOG_PATH

cp $COMPONENT_NAME $BIN_PATH/
chmod 700 $BIN_PATH/$COMPONENT_NAME

touch $CONFIG_FILE_PATH
chmod 600 $CONFIG_FILE_PATH

tenantctl config -v $env_file
if [ $? -ne 0 ]; then
  echo "Failed to update Tenant CLI configuration file"
  exit 1
fi
echo "Installation completed successfully!"

