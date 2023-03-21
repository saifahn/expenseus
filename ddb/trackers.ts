import { GetCommand, PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { CreateTrackerInput } from 'pages/api/v1/trackers';
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

export const makeTrackerIdKey = (id: string) => `${trackerKeyPrefix}#${id}`;

export type Tracker = {
  id: string;
  name: string;
  users: string[];
};

export type TrackerItem = {
  [tablePartitionKey]: string;
  [tableSortKey]: string;
  EntityType: typeof trackerEntityType;
  ID: string;
  Name: string;
  Users: string[];
  [gsi1PartitionKey]: typeof allTrackersKey;
  [gsi1SortKey]: string;
};

export function makeTrackerRepository({ ddb, tableName }: DDBWithConfig) {
  async function createTracker({ users, name }: CreateTrackerInput) {
    const id = ulid();
    const trackerIdKey = makeTrackerIdKey(id);
    for (const user of users) {
      const userIdKey = makeUserIdKey(user);

      await ddb.send(
        new PutCommand({
          TableName: tableName,
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

    await ddb.send(
      new PutCommand({
        TableName: tableName,
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

  async function getTracker(id: string) {
    const trackerIdKey = makeTrackerIdKey(id);

    const result = await ddb.send(
      new GetCommand({
        TableName: tableName,
        Key: {
          [tablePartitionKey]: trackerIdKey,
          [tableSortKey]: trackerIdKey,
        },
      }),
    );

    return result.Item as TrackerItem;
  }

  async function getTrackersByUser(userId: string) {
    const userIdKey = makeUserIdKey(userId);
    const allTrackerPrefix = `${trackerKeyPrefix}#`;

    const result = await ddb.send(
      new QueryCommand({
        TableName: tableName,
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

    return (result.Items as TrackerItem[]) ?? [];
  }

  return {
    createTracker,
    getTracker,
    getTrackersByUser,
  };
}
