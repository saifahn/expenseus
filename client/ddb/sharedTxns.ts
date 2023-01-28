import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';
import { DeleteCommand, PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { SubcategoryKey } from 'data/categories';
import { SharedTxn } from 'pages/shared/trackers/[trackerId]';
import { monotonicFactory } from 'ulid';
import { ItemDoesNotExistError } from './errors';
import {
  DDBWithConfig,
  gsi1PartitionKey,
  gsi1SortKey,
  gsi1Name,
  tablePartitionKey,
  tableSortKey,
  unsettledTxnsIndexName,
  unsettledTxnsIndexPK,
  unsettledTxnsIndexSK,
} from './schema';
import { makeTrackerIdKey } from './trackers';
import { makeUserIdKey } from './users';

const sharedTxnKeyPrefix = 'txn.shared',
  sharedTxnEntityType = 'sharedTransaction',
  unsettledFlagTrue = 'X';

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
  Unsettled?: typeof unsettledFlagTrue;
};

const ulid = monotonicFactory();

/**
 * Creates a shared transaction based on the given input.
 */
export async function createSharedTxn(d: DDBWithConfig, txn: SharedTxn) {
  // TODO: update the 2nd arg type so that it has to be without an ID
  const txnId = ulid(txn.date);
  const txnIdKey = makeSharedTxnIdKey(txnId);
  const txnDateIdKey = makeSharedTxnDateIdKey(txn);

  // TODO: add type for this as well?
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
    ...(txn.unsettled && { Unsettled: unsettledFlagTrue }),
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

/**
 * Updates a shared transaction based on the given input.
 */
export async function updateSharedTxn(d: DDBWithConfig, txn: SharedTxn) {
  // TODO: update the 2nd arg type so it has to have an ID
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
    ...(txn.unsettled && { Unsettled: unsettledFlagTrue }),
  };

  for (const user of txn.participants) {
    const userIdKey = makeUserIdKey(user);
    const uItem = {
      ...item,
      [tablePartitionKey]: userIdKey,
      [gsi1PartitionKey]: userIdKey,
    };
    try {
      await d.ddb.send(
        new PutCommand({
          TableName: d.tableName,
          Item: uItem,
          ExpressionAttributeNames: {
            '#SK': tableSortKey,
          },
          ConditionExpression: 'attribute_exists(#SK)',
        }),
      );
    } catch (err) {
      if (err instanceof ConditionalCheckFailedException) {
        throw new ItemDoesNotExistError('shared txn does not exist');
      }
      throw err;
    }
  }

  const trackerIdKey = makeTrackerIdKey(txn.tracker);
  const trackerItem = {
    ...item,
    [tablePartitionKey]: trackerIdKey,
    [gsi1PartitionKey]: trackerIdKey,
  };
  try {
    await d.ddb.send(
      new PutCommand({
        TableName: d.tableName,
        Item: trackerItem,
        ExpressionAttributeNames: {
          '#SK': tableSortKey,
        },
        ConditionExpression: 'attribute_exists(#SK)',
      }),
    );
  } catch (err) {
    if (err instanceof ConditionalCheckFailedException) {
      throw new ItemDoesNotExistError('shared txn does not exist');
    }
    throw err;
  }
}

type DeleteSharedTxnInput = {
  tracker: string;
  txnId: string;
  participants: string[];
};

/**
 * Deletes a shared transaction based on the given input
 */
export async function deleteSharedTxn(
  d: DDBWithConfig,
  input: DeleteSharedTxnInput,
) {
  const txnIdKey = makeSharedTxnIdKey(input.txnId);

  for (const user of input.participants) {
    const userIdKey = makeUserIdKey(user);
    await d.ddb.send(
      new DeleteCommand({
        TableName: d.tableName,
        Key: {
          [tablePartitionKey]: userIdKey,
          [tableSortKey]: txnIdKey,
        },
      }),
    );
  }

  const trackerIdKey = makeTrackerIdKey(input.tracker);
  await d.ddb.send(
    new DeleteCommand({
      TableName: d.tableName,
      Key: {
        [tablePartitionKey]: trackerIdKey,
        [tableSortKey]: txnIdKey,
      },
    }),
  );
}

/**
 * Retrieves all shared transactions from the tracker with the given tracker ID.
 */
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

export async function getUnsettledTxnsByTracker(
  d: DDBWithConfig,
  trackerId: string,
) {
  const trackerIdKey = makeTrackerIdKey(trackerId);

  const result = await d.ddb.send(
    new QueryCommand({
      TableName: d.tableName,
      IndexName: unsettledTxnsIndexName,
      ExpressionAttributeNames: {
        '#unsettledPK': unsettledTxnsIndexPK,
        '#unsettledSK': unsettledTxnsIndexSK,
      },
      ExpressionAttributeValues: {
        ':trackerIdKey': trackerIdKey,
        ':true': unsettledFlagTrue,
      },
      KeyConditionExpression:
        '#unsettledPK = :trackerIdKey and #unsettledSK = :true',
    }),
  );
  return (result.Items as SharedTxnItem[]) ?? [];
}
