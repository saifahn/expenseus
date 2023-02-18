import {
  CreateTableCommand,
  DeleteTableCommand,
  DynamoDBClient,
  ProjectionType,
  ResourceInUseException,
} from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient } from '@aws-sdk/lib-dynamodb';

// TODO: load from env variable
const endpoint = process.env.DDB_ENDPOINT;
const accessKeyId = process.env.DDB_ACCESS_KEY_ID ?? '';
const secretAccessKey = process.env.DDB_SECRET_ACCESS_KEY ?? '';
export const tableName = process.env.DDB_TABLE ?? 'expenseus-default';

export type DDBWithConfig = {
  ddb: DynamoDBDocumentClient;
  tableName: string;
};

export function setUpDdb(tableName: string) {
  const ddbClient = new DynamoDBClient({
    endpoint,
    credentials: {
      accessKeyId,
      secretAccessKey,
    },
  });
  const ddbDocClient = DynamoDBDocumentClient.from(ddbClient);
  return {
    ddb: ddbDocClient,
    tableName,
  };
}

export const tablePartitionKey = 'PK',
  tableSortKey = 'SK',
  gsi1Name = 'GSI1',
  gsi1PartitionKey = 'GSI1PK',
  gsi1SortKey = 'GSI1SK',
  unsettledTxnsIndexName = 'UnsettledTransactions',
  unsettledTxnsIndexPK = 'PK',
  unsettledTxnsIndexSK = 'Unsettled';

function createTableCommand(tableName: string) {
  return new CreateTableCommand({
    TableName: tableName,
    KeySchema: [
      {
        AttributeName: tablePartitionKey,
        KeyType: 'HASH',
      },
      {
        AttributeName: tableSortKey,
        KeyType: 'RANGE',
      },
    ],
    AttributeDefinitions: [
      {
        AttributeName: tablePartitionKey,
        AttributeType: 'S',
      },
      {
        AttributeName: tableSortKey,
        AttributeType: 'S',
      },
      {
        AttributeName: gsi1PartitionKey,
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
            AttributeName: gsi1PartitionKey,
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
  } catch (err) {
    if (err instanceof ResourceInUseException) {
      console.log(`Table exists. Continuing...`);
      return;
    }
    throw err;
  }
}

export async function TESTONLY_deleteTable(tableName: string) {
  const d = setUpDdb(tableName);
  try {
    await d.ddb.send(
      new DeleteTableCommand({
        TableName: tableName,
      }),
    );
  } catch (err) {
    if (err instanceof ResourceInUseException) {
      console.error('Tried to delete a non-existent table. Continuing...');
    }
    throw err;
  }
}
