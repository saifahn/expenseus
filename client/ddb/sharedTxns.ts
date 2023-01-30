import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';
import {
  DeleteCommand,
  PutCommand,
  QueryCommand,
  UpdateCommand,
} from '@aws-sdk/lib-dynamodb';
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

export function makeSharedTxnRepository({ ddb, tableName }: DDBWithConfig) {
  /**
   * Creates a shared transaction based on the given input.
   */
  async function createSharedTxn(txn: SharedTxn) {
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
      await ddb.send(
        new PutCommand({
          TableName: tableName,
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
    await ddb.send(
      new PutCommand({
        TableName: tableName,
        Item: trackerItem,
      }),
    );
  }

  /**
   * Updates a shared transaction based on the given input.
   */
  async function updateSharedTxn(txn: SharedTxn) {
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
        await ddb.send(
          new PutCommand({
            TableName: tableName,
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
      await ddb.send(
        new PutCommand({
          TableName: tableName,
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
  async function deleteSharedTxn(input: DeleteSharedTxnInput) {
    const txnIdKey = makeSharedTxnIdKey(input.txnId);

    for (const user of input.participants) {
      const userIdKey = makeUserIdKey(user);
      await ddb.send(
        new DeleteCommand({
          TableName: tableName,
          Key: {
            [tablePartitionKey]: userIdKey,
            [tableSortKey]: txnIdKey,
          },
        }),
      );
    }

    const trackerIdKey = makeTrackerIdKey(input.tracker);
    await ddb.send(
      new DeleteCommand({
        TableName: tableName,
        Key: {
          [tablePartitionKey]: trackerIdKey,
          [tableSortKey]: txnIdKey,
        },
      }),
    );
  }

  type SettleTxnInput = {
    id: string;
    trackerId: string;
    participants: string[];
  };

  /**
   * Settles transactions for the given input.
   */
  async function settleTxns(txns: SettleTxnInput[]) {
    for (const txn of txns) {
      const trackerIdKey = makeTrackerIdKey(txn.trackerId);
      const txnIdKey = makeSharedTxnIdKey(txn.id);
      await ddb.send(
        new UpdateCommand({
          TableName: tableName,
          Key: {
            [tablePartitionKey]: trackerIdKey,
            [tableSortKey]: txnIdKey,
          },
          ExpressionAttributeNames: {
            '#unsettled': 'Unsettled',
          },
          UpdateExpression: 'REMOVE #unsettled',
        }),
      );

      for (const user of txn.participants) {
        const userIdKey = makeUserIdKey(user);
        await ddb.send(
          new UpdateCommand({
            TableName: tableName,
            Key: {
              [tablePartitionKey]: userIdKey,
              [tableSortKey]: txnIdKey,
            },
            ExpressionAttributeNames: {
              '#unsettled': 'Unsettled',
            },
            UpdateExpression: 'REMOVE #unsettled',
          }),
        );
      }
    }
  }

  /**
   * Retrieves all shared transactions from the tracker with the given tracker I
   */
  async function getTxnsByTracker(trackerId: string) {
    const trackerIdKey = makeTrackerIdKey(trackerId);

    const result = await ddb.send(
      new QueryCommand({
        TableName: tableName,
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

  async function getUnsettledTxnsByTracker(trackerId: string) {
    const trackerIdKey = makeTrackerIdKey(trackerId);

    const result = await ddb.send(
      new QueryCommand({
        TableName: tableName,
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

  return {
    createSharedTxn,
    updateSharedTxn,
    deleteSharedTxn,
    settleTxns,
    getTxnsByTracker,
    getUnsettledTxnsByTracker,
  };
}
