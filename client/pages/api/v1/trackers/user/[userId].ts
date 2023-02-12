import { setUpTrackerRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch } from 'utils/withTryCatch';

export default async function getTrackersByUserHandler(
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
  if (session.user?.email !== req.query.userId) {
    return res
      .status(403)
      .json({ error: "you cannot view another user's trackers" });
  }

  const trackerRepo = setUpTrackerRepo();
  const [items, err] = await withAsyncTryCatch(
    trackerRepo.getTrackersByUser(session.user?.email!),
  );
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong while retrieving trackers' });
  }
  const trackers = items?.map((t) => ({
    id: t.ID,
    name: t.Name,
    users: t.Users,
  }));
  return res.status(200).json(trackers);
}
