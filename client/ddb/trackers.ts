import { GetCommand, PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { ulid } from 'ulid';
import {
  DDBWithConfig,
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from './schema';
import { makeUserIdKey } from './users';

const trackerKeyPrefix = 'tracker',
  trackerEntityType = 'tracker',
  allTrackersKey = 'trackers';
const makeTrackerIdKey = (id: string) => `${trackerKeyPrefix}#${id}`;

export async function createTracker(
  d: DDBWithConfig,
  { users, name }: { users: string[]; name: string },
) {
  const id = ulid();
  const trackerIdKey = makeTrackerIdKey(id);
  for (const user of users) {
    const userIdKey = makeUserIdKey(user);

    await d.ddb.send(
      new PutCommand({
        TableName: d.tableName,
        Item: {
          [tablePartitionKey]: userIdKey,
          [tableSortKey]: trackerIdKey,
          EntityType: trackerEntityType,
          ID: id,
          Name: name,
          Users: users,
          [gsi1PartitionKey]: allTrackersKey,
          [gsi1SortKey]: trackerIdKey,
        },
      }),
    );
  }

  await d.ddb.send(
    new PutCommand({
      TableName: d.tableName,
      Item: {
        [tablePartitionKey]: trackerIdKey,
        [tableSortKey]: trackerIdKey,
        EntityType: trackerEntityType,
        ID: id,
        Name: name,
        Users: users,
        [gsi1PartitionKey]: allTrackersKey,
        [gsi1SortKey]: trackerIdKey,
      },
    }),
  );
}

export async function getTracker(d: DDBWithConfig, id: string) {
  const trackerIdKey = makeTrackerIdKey(id);

  const result = await d.ddb.send(
    new GetCommand({
      TableName: d.tableName,
      Key: {
        [tablePartitionKey]: trackerIdKey,
        [tableSortKey]: trackerIdKey,
      },
    }),
  );

  return result.Item;
}

export async function getTrackersByUser(d: DDBWithConfig, userId: string) {
  const userIdKey = makeUserIdKey(userId);
  const allTrackerPrefix = `${trackerKeyPrefix}#`;

  const result = await d.ddb.send(
    new QueryCommand({
      TableName: d.tableName,
      ExpressionAttributeNames: {
        '#PK': tablePartitionKey,
        '#SK': tableSortKey,
      },
      ExpressionAttributeValues: {
        ':userKey': userIdKey,
        ':allTrackerPrefix': allTrackerPrefix,
      },
      KeyConditionExpression:
        '#PK = :userKey and begins_with(#SK, :allTrackerPrefix)',
    }),
  );

  return result.Items;
}
