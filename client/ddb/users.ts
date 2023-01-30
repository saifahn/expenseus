import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';
import { PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
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
    return res.Items!;
  }

  return { createUser, getAllUsers };
}
