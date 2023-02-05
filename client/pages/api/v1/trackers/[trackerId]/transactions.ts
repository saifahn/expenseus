import { sharedTxnItemToModel } from 'ddb/itemToModel';
import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch } from 'utils/withTryCatch';

export default async function txnsByTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (!['GET', 'POST'].includes(req.method!)) {
    return res.status(405).json({ error: 'invalid method' });
  }

  const session = await getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  const sharedTxnRepo = setUpSharedTxnRepo();
  const [items, err] = await withAsyncTryCatch(
    sharedTxnRepo.getTxnsByTracker(req.query.trackerId as string),
  );
  if (err) {
    return res
      .status(500)
      .json({
        error: 'something went wrong while getting transactions from tracker',
      });
  }
  const txns = items?.map(sharedTxnItemToModel);
  return res.status(200).json(txns);
}
