import { setUpDdb } from 'ddb/schema';
import { makeUserRepository, userItemsToUsers } from 'ddb/users';
import { NextApiRequest, NextApiResponse } from 'next';

export default async function usersHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    res.status(405).send('method not allowed');
    return;
  }
  // TODO: get database name from env
  const ddb = setUpDdb('test-ddb');
  const userRepo = makeUserRepository(ddb);
  const userItems = await userRepo.getAllUsers();
  const users = userItemsToUsers(userItems);

  res.status(200).json(users);
}
