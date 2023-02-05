import { setUpDdb, tableName } from './schema';
import { makeTxnRepository } from './txns';
import { makeUserRepository } from './users';

export const setUpTxnRepo = () => makeTxnRepository(setUpDdb(tableName));
export const setUpUserRepo = () => makeUserRepository(setUpDdb(tableName));
