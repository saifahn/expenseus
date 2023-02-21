import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';
import {
  DeleteCommand,
  PutCommand,
  QueryCommand,
  UpdateCommand,
} from '@aws-sdk/lib-dynamodb';
import { SubcategoryKey } from 'data/categories';
import { CreateSharedTxnPayload } from 'pages/api/v1/trackers/[trackerId]/transactions';
import { UpdateSharedTxnPayload } from 'pages/api/v1/trackers/[trackerId]/transactions/[txnId]';
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
const makeSharedTxnDateKey = (date: number) =>
  `${sharedTxnKeyPrefix}#${date.toString()}`;
const makeSharedTxnDateIdKey = ({ date, id }: { date: number; id: string }) =>
  `${sharedTxnKeyPrefix}#${date}#${id}`;

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
  SplitJSON?: string;
};

export type SharedTxn = {
  id: string;
  date: number;
  amount: number;
  location: string;
  tracker: string;
  category: SubcategoryKey;
  participants: string[];
  payer: string;
  details: string;
  unsettled?: boolean;
  split?: {
    [k: string]: number;
  };
};

const ulid = monotonicFactory();

export function makeSharedTxnRepository({ ddb, tableName }: DDBWithConfig) {
  /**
   * Creates a shared transaction based on the given input.
   */
  async function createSharedTxn(txn: CreateSharedTxnPayload) {
    const txnId = ulid(txn.date);
    const txnIdKey = makeSharedTxnIdKey(txnId);
    const txnDateIdKey = makeSharedTxnDateIdKey({ id: txnId, date: txn.date });

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
      // storing a JS object using the ddb Map type would be better, but the original
      // implementation used a string so this is kept for compatibility
      ...(txn.split && { SplitJSON: JSON.stringify(txn.split) }),
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
  async function updateSharedTxn(txn: UpdateSharedTxnPayload) {
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
      ...(txn.split && { SplitJSON: JSON.stringify(txn.split) }),
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

  /**
   * Retrieves all shared transactions from the tracker with the given tracker ID.
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

  type ByTrackerBetweenDatesInput = {
    tracker: string;
    from: number;
    to: number;
  };
  /**
   * Retrieves transactions between the given dates for the given tracker ID.
   */
  async function getTxnsByTrackerBetweenDates({
    tracker,
    from,
    to,
  }: ByTrackerBetweenDatesInput) {
    const trackerIdKey = makeTrackerIdKey(tracker);
    const txnDateFromKey = makeSharedTxnDateKey(from);
    const txnDateToKey = makeSharedTxnDateKey(to + 1); // to be inclusive for the date

    const results = await ddb.send(
      new QueryCommand({
        TableName: tableName,
        IndexName: gsi1Name,
        ExpressionAttributeNames: {
          '#GSI1PK': gsi1PartitionKey,
          '#GSI1SK': gsi1SortKey,
        },
        ExpressionAttributeValues: {
          ':trackerKey': trackerIdKey,
          ':txnDateFromKey': txnDateFromKey,
          ':txnDateToKey': txnDateToKey,
        },
        KeyConditionExpression:
          '#GSI1PK = :trackerKey and #GSI1SK BETWEEN :txnDateFromKey AND :txnDateToKey',
      }),
    );

    return (results.Items as SharedTxnItem[]) ?? [];
  }

  type ByUserBetweenDatesInput = {
    user: string;
    from: number;
    to: number;
  };
  /**
   * Retrieves shared transactions between the given dates for the given user ID.
   */
  async function getSharedTxnsByUserBetweenDates({
    user,
    from,
    to,
  }: ByUserBetweenDatesInput) {
    const userIdKey = makeUserIdKey(user);
    const txnDateFromKey = makeSharedTxnDateKey(from);
    const txnDateToKey = makeSharedTxnDateKey(to);

    const results = await ddb.send(
      new QueryCommand({
        TableName: tableName,
        IndexName: gsi1Name,
        ExpressionAttributeNames: {
          '#GSI1PK': gsi1PartitionKey,
          '#GSI1SK': gsi1SortKey,
        },
        ExpressionAttributeValues: {
          ':trackerKey': userIdKey,
          ':txnDateFromKey': txnDateFromKey,
          ':txnDateToKey': txnDateToKey,
        },
        KeyConditionExpression:
          '#GSI1PK = :trackerKey and #GSI1SK BETWEEN :txnDateFromKey AND :txnDateToKey',
      }),
    );

    return (results.Items as SharedTxnItem[]) ?? [];
  }

  /**
   * Retrieves unsettled transactions from the tracker with the given tracker ID.
   */
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

  return {
    createSharedTxn,
    updateSharedTxn,
    deleteSharedTxn,
    getTxnsByTracker,
    getTxnsByTrackerBetweenDates,
    getSharedTxnsByUserBetweenDates,
    getUnsettledTxnsByTracker,
    settleTxns,
  };
}
