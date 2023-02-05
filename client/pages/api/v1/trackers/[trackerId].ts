import { setUpTrackerRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch } from 'utils/withTryCatch';

export default async function getTrackerByIdHandler(
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

  const trackerRepo = setUpTrackerRepo();
  const [item, err] = await withAsyncTryCatch(
    trackerRepo.getTracker(req.query.trackerId as string),
  );
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong while retrieving tracker' });
  }
  if (!item) {
    return res.status(404).json({ error: 'no tracker with that ID was found' });
  }
}
