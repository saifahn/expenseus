import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';
import { GetCommand, PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { User } from 'components/UserList';
import {
  DDBWithConfig,
  gsi1Name,
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from './schema';
import { UserAlreadyExistsError } from './errors';

const userKeyPrefix = 'user',
  userEntityType = 'user',
  allUsersKey = 'users';

export const makeUserIdKey = (id: string) => `${userKeyPrefix}#${id}`;

export type UserItem = {
  [tablePartitionKey]: string;
  [tableSortKey]: string;
  EntityType: typeof userEntityType;
  ID: string;
  Username: string;
  Name: string;
  [gsi1PartitionKey]: typeof allUsersKey;
  [gsi1SortKey]: string;
};

export function userItemsToUsers(items: UserItem[]): User[] {
  return items.map((i) => ({
    id: i.ID,
    username: i.Username,
    name: i.Name,
  }));
}

export function makeUserRepository({ ddb, tableName }: DDBWithConfig) {
  async function createUser(user: User) {
    const userIdKey = makeUserIdKey(user.id);
    const userItem = {
      [tablePartitionKey]: userIdKey,
      [tableSortKey]: userIdKey,
      EntityType: userEntityType,
      ID: user.id,
      Username: user.username,
      Name: user.name,
      [gsi1PartitionKey]: allUsersKey,
      [gsi1SortKey]: userIdKey,
    };
    try {
      const result = await ddb.send(
        new PutCommand({
          TableName: tableName,
          Item: userItem,
          ConditionExpression: 'attribute_not_exists(PK)',
        }),
      );
      return result;
    } catch (err) {
      if (err instanceof ConditionalCheckFailedException) {
        throw new UserAlreadyExistsError();
      }
      throw err;
    }
  }

  async function getUser(id: string) {
    const userIdKey = makeUserIdKey(id);
    try {
      const result = await ddb.send(
        new GetCommand({
          TableName: tableName,
          Key: {
            [tablePartitionKey]: userIdKey,
            [tableSortKey]: userIdKey,
          },
        }),
      );
      return result.Item;
    } catch (err) {
      console.error(err);
      throw err;
    }
  }

  async function getAllUsers() {
    const res = await ddb.send(
      new QueryCommand({
        TableName: tableName,
        IndexName: gsi1Name,
        ExpressionAttributeNames: {
          '#GSI1PK': gsi1PartitionKey,
        },
        ExpressionAttributeValues: {
          ':usersKey': allUsersKey,
        },
        KeyConditionExpression: '#GSI1PK = :usersKey',
      }),
    );
    return (res.Items! as UserItem[]) || [];
  }

  return { createUser, getUser, getAllUsers };
}
