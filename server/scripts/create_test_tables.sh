#!/bin/bash

SCRIPTS="$(dirname $0)"
ROOT="$(dirname "$SCRIPTS")"

# source the env
source $ROOT/.envrc

# create the transaction table
aws dynamodb create-table \
  --endpoint-url $DYNAMODB_ENDPOINT_LOCAL \
  --table-name $DYNAMODB_TRANSACTIONS_TABLE_NAME \
  --attribute-definitions \
  AttributeName=id,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

# create the users table
aws dynamodb create-table \
  --endpoint-url $DYNAMODB_ENDPOINT_LOCAL \
  --table-name $DYNAMODB_USERS_TABLE_NAME \
  --attribute-definitions \
  AttributeName=id,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
