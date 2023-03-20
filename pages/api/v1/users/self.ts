import { userItemToUser } from 'ddb/itemToModel';
import { setUpUserRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch } from 'utils/withTryCatch';
import { authOptions } from '../../auth/[...nextauth]';

export default async function getSelfHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const session = await getServerSession(req, res, authOptions);
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  const userRepo = setUpUserRepo();
  const [item, err] = await withAsyncTryCatch(
    userRepo.getUser(session.user?.email!),
  );
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong while fetching self' });
  }
  return res.status(200).json(userItemToUser(item!));
}
