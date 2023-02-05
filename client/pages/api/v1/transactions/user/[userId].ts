import { setUpTxnRepo } from 'ddb/setUpRepos';
import { txnItemToTxn } from 'ddb/itemToModel';
import { NextApiRequest, NextApiResponse } from 'next';
import { withAsyncTryCatch } from 'utils/withTryCatch';
import { getServerSession } from 'next-auth';

export default async function txnByUserIdHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const session = await getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  const sessionUser = session!.user?.email;
  const userId = req.query.userId as string;
  if (sessionUser !== userId) {
    return res
      .status(403)
      .json({ error: "you cannot retrieve other users' transactions" });
  }

  const txnRepo = setUpTxnRepo();
  const [items, err] = await withAsyncTryCatch(txnRepo.getTxnsByUserId(userId));
  if (err) {
    return res.status(500).json({
      error: 'something went wrong while trying to get transactions by user ID',
    });
  }
  const txns = items?.map(txnItemToTxn);
  return res.status(200).json(txns);
}
