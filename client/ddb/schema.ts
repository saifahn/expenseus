import {
  CreateTableCommand,
  DeleteTableCommand,
  DynamoDBClient,
  ResourceInUseException,
} from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient } from '@aws-sdk/lib-dynamodb';

// TODO: load from env variable
const endpoint = 'http://localhost:8000';

export function setUpDdb(tableName: string) {
  const ddbClient = new DynamoDBClient({ endpoint });
  const ddbDocClient = DynamoDBDocumentClient.from(ddbClient);
  return {
    ddb: ddbDocClient,
    tableName,
  };
}

function createTableCommand(tableName: string) {
  return new CreateTableCommand({
    TableName: tableName,
    KeySchema: [
      {
        AttributeName: 'PK',
        KeyType: 'HASH',
      },
      {
        AttributeName: 'SK',
        KeyType: 'RANGE',
      },
    ],
    AttributeDefinitions: [
      {
        AttributeName: 'PK',
        AttributeType: 'S',
      },
      {
        AttributeName: 'SK',
        AttributeType: 'S',
      },
      // {
      //   AttributeName: 'GSI1PK',
      //   AttributeType: 'S',
      // },
      // {
      //   AttributeName: 'GSI1SK',
      //   AttributeType: 'S',
      // },
      // {
      //   AttributeName: 'Unsettled',
      //   AttributeType: 'S',
      // },
    ],
    ProvisionedThroughput: {
      ReadCapacityUnits: 1,
      WriteCapacityUnits: 1,
    },
  });
}

export async function createTableIfNotExists(tableName: string) {
  const d = setUpDdb(tableName);
  try {
    await d.ddb.send(createTableCommand(tableName));
    console.log('Table created successfully.');
  } catch (err) {
    if (err instanceof ResourceInUseException) {
      console.log(`Table exists. Continuing...`);
      return;
    }
    throw err;
  }
}

export async function deleteTable(tableName: string) {
  const d = setUpDdb(tableName);
  try {
    await d.ddb.send(
      new DeleteTableCommand({
        TableName: tableName,
      }),
    );
    console.log('Table deleted.');
  } catch (err) {
    console.error('something went wrong while trying to delete the table');
  }
}
