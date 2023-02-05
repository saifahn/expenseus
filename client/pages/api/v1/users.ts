import { setUpUserRepo } from 'ddb/setUpRepos';
import { userItemsToUsers } from 'ddb/users';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';

export default async function usersHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const session = await getServerSession();

  if (!session) {
    res.status(401).json({ error: 'you need to be logged in' });
    return;
  }

  if (req.method !== 'GET') {
    res.status(405).json({ error: 'method not allowed' });
    return;
  }
  const userRepo = setUpUserRepo();
  const userItems = await userRepo.getAllUsers();
  const users = userItemsToUsers(userItems);

  res.status(200).json(users);
}
