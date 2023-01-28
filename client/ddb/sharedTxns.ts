import { PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { SubcategoryKey } from 'data/categories';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { monotonicFactory } from 'ulid';
import {
  DDBWithConfig,
  gsi1PartitionKey,
  gsi1SortKey,
  gsi1Name,
  tablePartitionKey,
  tableSortKey,
} from './schema';
import { makeTrackerIdKey } from './trackers';
import { makeUserIdKey } from './users';

const sharedTxnKeyPrefix = 'txn.shared',
  sharedTxnEntityType = 'sharedTransaction';

const makeSharedTxnIdKey = (id: string) => `${sharedTxnKeyPrefix}#${id}`;
const makeSharedTxnDateIdKey = (txn: SharedTxn) =>
  `${sharedTxnKeyPrefix}#${txn.date}#${txn.id}`;

export type SharedTxnItem = {
  [tablePartitionKey]: string;
  [tableSortKey]: string;
  [gsi1PartitionKey]: string;
  [gsi1SortKey]: string;
  EntityType: typeof sharedTxnEntityType;
  ID: string;
  Date: number;
  Amount: number;
  Location: string;
  Tracker: string;
  Category: SubcategoryKey;
  Participants: string[];
  Payer: string;
  Details: string;
};

const ulid = monotonicFactory();

export async function createSharedTxn(d: DDBWithConfig, txn: SharedTxn) {
  const txnId = ulid(txn.date);
  const txnIdKey = makeSharedTxnIdKey(txnId);
  const txnDateIdKey = makeSharedTxnDateIdKey(txn);

  const item = {
    [tableSortKey]: txnIdKey,
    [gsi1SortKey]: txnDateIdKey,
    EntityType: sharedTxnEntityType,
    ID: txnId,
    Category: txn.category,
    Tracker: txn.tracker,
    Participants: txn.participants,
    Date: txn.date,
    Amount: txn.amount,
    Location: txn.location,
    Payer: txn.payer,
    Details: txn.details,
  };

  // store the representation of the tracker under each user so trackers can
  // be retrieved for them
  for (const user of txn.participants) {
    const userIdKey = makeUserIdKey(user);
    const uItem = {
      ...item,
      [tablePartitionKey]: userIdKey,
      [gsi1PartitionKey]: userIdKey,
    };
    await d.ddb.send(
      new PutCommand({
        TableName: d.tableName,
        Item: uItem,
      }),
    );
  }

  const trackerIdKey = makeTrackerIdKey(txn.tracker);
  const trackerItem = {
    ...item,
    [tablePartitionKey]: trackerIdKey,
    [gsi1PartitionKey]: trackerIdKey,
  };
  await d.ddb.send(
    new PutCommand({
      TableName: d.tableName,
      Item: trackerItem,
    }),
  );
}

export async function updateSharedTxn(d: DDBWithConfig, txn: SharedTxn) {
  const txnIdKey = makeSharedTxnIdKey(txn.id);
  const txnDateIdKey = makeSharedTxnDateIdKey(txn);

  const item = {
    [tableSortKey]: txnIdKey,
    [gsi1SortKey]: txnDateIdKey,
    EntityType: sharedTxnEntityType,
    ID: txn.id,
    Category: txn.category,
    Tracker: txn.tracker,
    Participants: txn.participants,
    Date: txn.date,
    Amount: txn.amount,
    Location: txn.location,
    Payer: txn.payer,
    Details: txn.details,
  };

  for (const user of txn.participants) {
    const userIdKey = makeUserIdKey(user);
    const uItem = {
      ...item,
      [tablePartitionKey]: userIdKey,
      [gsi1PartitionKey]: userIdKey,
    };
    await d.ddb.send(
      new PutCommand({
        TableName: d.tableName,
        Item: uItem,
      }),
    );
  }

  const trackerIdKey = makeTrackerIdKey(txn.tracker);
  const trackerItem = {
    ...item,
    [tablePartitionKey]: trackerIdKey,
    [gsi1PartitionKey]: trackerIdKey,
  };
  await d.ddb.send(
    new PutCommand({
      TableName: d.tableName,
      Item: trackerItem,
    }),
  );
}

export async function getTxnsByTracker(d: DDBWithConfig, trackerId: string) {
  const trackerIdKey = makeTrackerIdKey(trackerId);

  const result = await d.ddb.send(
    new QueryCommand({
      TableName: d.tableName,
      IndexName: gsi1Name,
      ExpressionAttributeNames: {
        '#GSI1PK': gsi1PartitionKey,
        '#GSI1SK': gsi1SortKey,
      },
      ExpressionAttributeValues: {
        ':trackerIdKey': trackerIdKey,
        ':sharedTxnKeyPrefix': sharedTxnKeyPrefix,
      },
      // return in descending order
      ScanIndexForward: false,
      KeyConditionExpression:
        '#GSI1PK = :trackerIdKey and begins_with(#GSI1SK, :sharedTxnKeyPrefix)',
    }),
  );

  return (result.Items as SharedTxnItem[]) ?? [];
}
