import { userItemToUser } from 'ddb/itemToModel';
import { setUpUserRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { authOptions } from '../auth/[...nextauth]';

export default async function usersHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const session = await getServerSession(req, res, authOptions);

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
  const users = userItems.map(userItemToUser);

  res.status(200).json(users);
}
