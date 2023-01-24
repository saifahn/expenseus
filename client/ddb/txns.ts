import { Transaction } from 'types/Transaction';
import {
  ddbWithConfig,
  gsi1Name,
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from './schema';
import { monotonicFactory } from 'ulid';
import { makeUserIdKey } from './users';
import { SubcategoryKey } from 'data/categories';
import { PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';

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

export async function createTxn(d: ddbWithConfig, txn: Transaction) {
  const txnId = ulid(txn.date);
  const userIdKey = makeUserIdKey(txn.userId);
  const txnIdKey = makeTxnIdKey(txnId);
  const txnDateIdKey = makeTxnDateIdKey(txn);

  const txnItem: TxnItem = {
    [tablePartitionKey]: userIdKey,
    [tableSortKey]: txnIdKey,
    EntityType: txnEntityType,
    ID: txnId,
    UserID: txn.userId,
    Location: txn.location,
    Details: txn.details,
    Amount: txn.amount,
    Date: txn.date,
    [gsi1PartitionKey]: userIdKey,
    [gsi1SortKey]: txnDateIdKey,
    Category: txn.category,
  };

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

export async function getTxnsByUserId(d: ddbWithConfig, id: string) {
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
  return result.Items;
}
