import { createTableIfNotExists, tableName } from './schema';

createTableIfNotExists(tableName)
  .then(() => console.log('table successfully created'))
  .catch((err) => console.error('something went wrong', err));
