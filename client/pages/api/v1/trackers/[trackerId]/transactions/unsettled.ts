import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch } from 'utils/withTryCatch';

export default async function getUnsettledTxnsByTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const session = getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  const sharedTxnRepo = setUpSharedTxnRepo();
  const [items, err] = await withAsyncTryCatch(
    sharedTxnRepo.getUnsettledTxnsByTracker(req.query.trackerId as string),
  );
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong while getting unsettled txns' });
  }
}
