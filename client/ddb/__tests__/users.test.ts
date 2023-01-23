import { createTableIfNotExists, deleteTable, setUpDdb } from 'ddb/schema';
import { createUser, getAllUsers } from 'ddb/users';

const d = setUpDdb('test-table');

describe('Users', () => {
  beforeEach(async () => {
    await createTableIfNotExists('test-table');
  });

  afterEach(async () => {
    await deleteTable('test-table');
  });

  it('creates a user correctly', async () => {
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
});
