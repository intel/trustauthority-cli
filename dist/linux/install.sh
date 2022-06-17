#!/bin/bash

COMPONENT_NAME=tenantctl
SERVICE_USERNAME=tac
SERVICE_ENV=tac.env

if [[ $EUID -ne 0 ]]; then
    echo "This installer must be run as root"
    exit 1
fi

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

echo "Setting up Tenant CLI Linux User..."
# useradd -M -> this user has no home directory
id -u $SERVICE_USERNAME 2> /dev/null || useradd -M --system --shell /sbin/nologin $SERVICE_USERNAME

echo "Installing Tenant CLI..."

PRODUCT_HOME=/opt/$SERVICE_USERNAME
BIN_PATH=$PRODUCT_HOME/bin
CONFIG_PATH=/etc/$SERVICE_USERNAME/

for directory in $BIN_PATH $CONFIG_PATH; do
  # mkdir -p will return 0 if directory exists or is a symlink to an existing directory or directory and parents can be created
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo "Cannot create directory: $directory"
    exit 1
  fi
  chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $directory
  chmod 700 $directory
done

cp $COMPONENT_NAME $BIN_PATH/ && chown $SERVICE_USERNAME:$SERVICE_USERNAME $BIN_PATH/*
chmod 700 $BIN_PATH/*
ln -sfT $BIN_PATH/$COMPONENT_NAME /usr/bin/$COMPONENT_NAME

chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME

tenantctl config -v $env_file
if [ $? -ne 0 ]; then
  echo "Failed to update Tenant CLI configuration file"
  exit 1
fi
echo "Installation completed successfully!"

