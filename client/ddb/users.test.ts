import { createTableIfNotExists, deleteTable, setUpDdb } from 'ddb/schema';
import { createUser, getAllUsers, UserAlreadyExistsError } from 'ddb/users';

const userTestTable = 'user-test-table';
const d = setUpDdb(userTestTable);

describe('Users', () => {
  beforeEach(async () => {
    await createTableIfNotExists(userTestTable);
  });

  afterEach(async () => {
    await deleteTable(userTestTable);
  });

  test('a user can be created correctly', async () => {
    let users = await getAllUsers(d);
    expect(users).toHaveLength(0);

    const testUser = {
      id: 'test-user',
      username: 'TestUser',
      name: 'Test User',
    };
    await createUser(d, testUser);

    users = await getAllUsers(d);
    expect(users).toHaveLength(1);
    const expected = {
      EntityType: 'user',
      GSI1PK: 'users',
      Username: 'TestUser',
      GSI1SK: 'user#test-user',
      SK: 'user#test-user',
      PK: 'user#test-user',
      ID: 'test-user',
      Name: 'Test User',
    };
    expect(users).toContainEqual(expected);
  });

  test('a UserAlreadyExistsError is thrown when trying to create a user with an existing id', async () => {
    const testUser = {
      id: 'test-user',
      username: 'TestUser',
      name: 'Test User',
    };
    await createUser(d, testUser);

    const users = await getAllUsers(d);
    expect(users).toHaveLength(1);

    expect(
      createUser(d, {
        id: 'test-user',
        username: 'different-name',
        name: 'Someone Different?',
      }),
    ).rejects.toThrowError(UserAlreadyExistsError);
  });
});