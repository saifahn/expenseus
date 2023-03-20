import { setUpDdb, tableName } from './schema';
import { makeSharedTxnRepository } from './sharedTxns';
import { makeTrackerRepository } from './trackers';
import { makeTxnRepository } from './txns';
import { makeUserRepository } from './users';

export const setUpTxnRepo = () => makeTxnRepository(setUpDdb(tableName));
export const setUpUserRepo = () => makeUserRepository(setUpDdb(tableName));
export const setUpSharedTxnRepo = () =>
  makeSharedTxnRepository(setUpDdb(tableName));
export const setUpTrackerRepo = () =>
  makeTrackerRepository(setUpDdb(tableName));
