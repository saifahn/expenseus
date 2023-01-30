import { Transaction } from 'types/Transaction';
import {
  DDBWithConfig,
  gsi1Name,
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from './schema';
import { monotonicFactory } from 'ulid';
import { makeUserIdKey } from './users';
import { SubcategoryKey } from 'data/categories';
import {
  DeleteCommand,
  GetCommand,
  PutCommand,
  QueryCommand,
} from '@aws-sdk/lib-dynamodb';
import { ItemDoesNotExistError } from './errors';
import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';

const txnKeyPrefix = 'txn',
  txnEntityType = 'transaction',
  allTxnKey = 'transactions';

const makeTxnIdKey = (id: string) => `${txnKeyPrefix}#${id}`;
const makeTxnDateIdKey = (txn: Transaction) =>
  `${txnKeyPrefix}#${txn.date}#${txn.id}`;
const makeTxnDateKey = (date: number) => `${txnKeyPrefix}#${date}`;

export type TxnItem = {
  [tablePartitionKey]: string;
  [tableSortKey]: string;
  [gsi1PartitionKey]: string;
  [gsi1SortKey]: string;
  EntityType: typeof txnEntityType;
  ID: string;
  UserID: string;
  Date: number;
  Amount: number;
  Location: string;
  Category: SubcategoryKey;
  Details: string;
};

function txnToTxnItem(txn: Transaction): TxnItem {
  const userIdKey = makeUserIdKey(txn.userId);
  const txnIdKey = makeTxnIdKey(txn.id);
  const txnDateIdKey = makeTxnDateIdKey(txn);

  return {
    [tablePartitionKey]: userIdKey,
    [tableSortKey]: txnIdKey,
    EntityType: txnEntityType,
    ID: txn.id,
    UserID: txn.userId,
    Location: txn.location,
    Details: txn.details,
    Amount: txn.amount,
    Date: txn.date,
    [gsi1PartitionKey]: userIdKey,
    [gsi1SortKey]: txnDateIdKey,
    Category: txn.category,
  };
}

const ulid = monotonicFactory();

export function makeTxnRepository({ ddb, tableName }: DDBWithConfig) {
  async function createTxn(txn: Transaction) {
    const txnId = ulid(txn.date);
    txn.id = txnId;
    const txnItem = txnToTxnItem(txn);

    await ddb.send(
      new PutCommand({
        TableName: tableName,
        Item: txnItem,
        ExpressionAttributeNames: {
          '#SK': tableSortKey,
        },
        ConditionExpression: 'attribute_not_exists(#SK)',
      }),
    );
  }

  async function getTxn({ txnId, userId }: { txnId: string; userId: string }) {
    const userIdKey = makeUserIdKey(userId);
    const txnIdKey = makeTxnIdKey(txnId);

    const results = await ddb.send(
      new GetCommand({
        TableName: tableName,
        Key: {
          [tablePartitionKey]: userIdKey,
          [tableSortKey]: txnIdKey,
        },
      }),
    );
    return results.Item as TxnItem;
  }

  async function updateTxn(txn: Transaction) {
    const txnItem = txnToTxnItem(txn);

    try {
      await ddb.send(
        new PutCommand({
          TableName: tableName,
          Item: txnItem,
          ExpressionAttributeNames: {
            '#SK': tableSortKey,
          },
          ConditionExpression: 'attribute_exists(#SK)',
        }),
      );
    } catch (err) {
      if (err instanceof ConditionalCheckFailedException) {
        throw new ItemDoesNotExistError();
      }
      throw err;
    }
  }

  async function deleteTxn({
    txnId,
    userId,
  }: {
    txnId: string;
    userId: string;
  }) {
    const userIdKey = makeUserIdKey(userId);
    const txnIdKey = makeTxnIdKey(txnId);
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

  async function getTxnsByUserId(id: string) {
    const userIdKey = makeUserIdKey(id);
    const allTxnsPrefix = `${txnKeyPrefix}#`;

    const result = await ddb.send(
      new QueryCommand({
        TableName: tableName,
        IndexName: gsi1Name,
        ExpressionAttributeNames: {
          '#GSI1PK': gsi1PartitionKey,
          '#GSI1SK': gsi1SortKey,
        },
        ExpressionAttributeValues: {
          ':userKey': userIdKey,
          ':allTxnPrefix': allTxnsPrefix,
        },
        KeyConditionExpression:
          '#GSI1PK = :userKey and begins_with(#GSI1SK, :allTxnPrefix)',
        // return in descending order
        ScanIndexForward: false,
      }),
    );
    return (result.Items as TxnItem[]) ?? [];
  }

  async function getBetweenDates({
    userId,
    from,
    to,
  }: {
    userId: string;
    from: number;
    to: number;
  }) {
    const userIdKey = makeUserIdKey(userId);
    const txnDateFromKey = makeTxnDateKey(from);
    const txnDateToKey = makeTxnDateKey(to);

    const result = await ddb.send(
      new QueryCommand({
        TableName: tableName,
        IndexName: gsi1Name,
        ExpressionAttributeNames: {
          '#GSI1PK': gsi1PartitionKey,
          '#GSI1SK': gsi1SortKey,
        },
        ExpressionAttributeValues: {
          ':userKey': userIdKey,
          ':txnDateFromKey': txnDateFromKey,
          ':txnDateToKey': txnDateToKey,
        },
        KeyConditionExpression:
          '#GSI1PK = :userKey and #GSI1SK BETWEEN :txnDateFromKey AND :txnDateToKey',
      }),
    );
    return (result.Items as TxnItem[]) || [];
  }

  return {
    createTxn,
    getTxn,
    updateTxn,
    deleteTxn,
    getTxnsByUserId,
    getBetweenDates,
  };
}
