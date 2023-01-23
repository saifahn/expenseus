import { createTableIfNotExists, deleteTable } from 'ddb/schema';

describe('Users', () => {
  afterEach(async () => {
    await deleteTable('test-table');
  });

  it('creates a table successfully', async () => {
    expect(createTableIfNotExists('test-table')).resolves.not.toThrowError();
  });
});
