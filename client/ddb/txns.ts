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
import { DeleteCommand, PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { ItemDoesNotExistError } from './errors';
import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';

const txnKeyPrefix = 'txn',
  txnEntityType = 'transaction',
  allTxnKey = 'transactions';

const makeTxnIdKey = (id: string) => `${txnKeyPrefix}#${id}`;
const makeTxnDateIdKey = (txn: Transaction) =>
  `${txnKeyPrefix}#${txn.date}#${txn.id}`;

export type TxnItem = {
  [tablePartitionKey]: string;
  [tableSortKey]: string;
  EntityType: typeof txnEntityType;
  ID: string;
  UserID: string;
  Location: string;
  Details: string;
  Amount: number;
  Date: number;
  [gsi1PartitionKey]: string;
  [gsi1SortKey]: string;
  Category: SubcategoryKey;
};

const ulid = monotonicFactory();

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

export async function createTxn(d: DDBWithConfig, txn: Transaction) {
  const txnId = ulid(txn.date);
  txn.id = txnId;
  const txnItem = txnToTxnItem(txn);

  await d.ddb.send(
    new PutCommand({
      TableName: d.tableName,
      Item: txnItem,
      ExpressionAttributeNames: {
        '#SK': tableSortKey,
      },
      ConditionExpression: 'attribute_not_exists(#SK)',
    }),
  );
}

export async function updateTxn(d: DDBWithConfig, txn: Transaction) {
  const txnItem = txnToTxnItem(txn);

  try {
    await d.ddb.send(
      new PutCommand({
        TableName: d.tableName,
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

export async function deleteTxn(
  d: DDBWithConfig,
  { txnId, userId }: { txnId: string; userId: string },
) {
  const userIdKey = makeUserIdKey(userId);
  const txnIdKey = makeTxnIdKey(txnId);
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

export async function getTxnsByUserId(d: DDBWithConfig, id: string) {
  const userIdKey = makeUserIdKey(id);
  const allTxnsPrefix = `${txnKeyPrefix}#`;

  const result = await d.ddb.send(
    new QueryCommand({
      TableName: d.tableName,
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
