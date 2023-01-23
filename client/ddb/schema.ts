import {
  CreateTableCommand,
  DeleteTableCommand,
  DynamoDBClient,
  ProjectionType,
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

const tablePrimaryKey = 'PK',
  tableSortKey = 'SK',
  gsi1Name = 'GSI1',
  gsi1PrimaryKey = 'GSI1PK',
  gsi1SortKey = 'GSI1SK',
  unsettledTxnsIndexName = 'UnsettledTransactions',
  unsettledTxnsIndexPK = 'PK',
  unsettledTxnsIndexSK = 'Unsettled';

function createTableCommand(tableName: string) {
  return new CreateTableCommand({
    TableName: tableName,
    KeySchema: [
      {
        AttributeName: tablePrimaryKey,
        KeyType: 'HASH',
      },
      {
        AttributeName: tableSortKey,
        KeyType: 'RANGE',
      },
    ],
    AttributeDefinitions: [
      {
        AttributeName: tablePrimaryKey,
        AttributeType: 'S',
      },
      {
        AttributeName: tableSortKey,
        AttributeType: 'S',
      },
      {
        AttributeName: gsi1PrimaryKey,
        AttributeType: 'S',
      },
      {
        AttributeName: gsi1SortKey,
        AttributeType: 'S',
      },
      {
        AttributeName: unsettledTxnsIndexSK,
        AttributeType: 'S',
      },
    ],
    GlobalSecondaryIndexes: [
      {
        IndexName: gsi1Name,
        KeySchema: [
          {
            AttributeName: gsi1PrimaryKey,
            KeyType: 'HASH',
          },
          {
            AttributeName: gsi1SortKey,
            KeyType: 'RANGE',
          },
        ],
        Projection: {
          ProjectionType: ProjectionType.ALL,
        },
        ProvisionedThroughput: {
          ReadCapacityUnits: 1,
          WriteCapacityUnits: 1,
        },
      },
      {
        IndexName: unsettledTxnsIndexName,
        KeySchema: [
          {
            AttributeName: unsettledTxnsIndexPK,
            KeyType: 'HASH',
          },
          {
            AttributeName: unsettledTxnsIndexSK,
            KeyType: 'RANGE',
          },
        ],
        Projection: {
          ProjectionType: ProjectionType.ALL,
        },
        ProvisionedThroughput: {
          ReadCapacityUnits: 1,
          WriteCapacityUnits: 1,
        },
      },
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
