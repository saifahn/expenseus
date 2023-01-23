import { createTableIfNotExists, deleteTable } from 'ddb/schema';

describe('Users', () => {
  beforeEach(async () => {
    await createTableIfNotExists('test-table');
  });

  afterEach(async () => {
    await deleteTable('test-table');
  });

  it('creates and tears down properly', async () => {});
});
