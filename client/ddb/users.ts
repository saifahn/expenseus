import { ConditionalCheckFailedException } from '@aws-sdk/client-dynamodb';
import { PutCommand, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { User } from 'components/UserList';
import {
  ddbWithConfig,
  gsi1Name,
  gsi1PartitionKey,
  gsi1SortKey,
  tablePartitionKey,
  tableSortKey,
} from './schema';

const userKeyPrefix = 'user',
  userEntityType = 'user',
  allUsersKey = 'users';

export const makeUserIdKey = (id: string) => `${userKeyPrefix}#${id}`;

export class UserAlreadyExistsError extends Error {
  constructor() {
    super();
  }
}

export async function createUser(d: ddbWithConfig, user: User) {
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
    const result = await d.ddb.send(
      new PutCommand({
        TableName: d.tableName,
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

export async function getAllUsers(d: ddbWithConfig) {
  const res = await d.ddb.send(
    new QueryCommand({
      TableName: d.tableName,
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
