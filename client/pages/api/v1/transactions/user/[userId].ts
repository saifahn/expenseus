import { setUpTxnRepo } from 'ddb/setUpRepos';
import { txnItemToTxn } from 'ddb/itemToModel';
import { NextApiRequest, NextApiResponse } from 'next';
import { withAsyncTryCatch } from 'utils/withTryCatch';

export default async function txnByUserIdHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const userId = req.query.userId as string;
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
