import { setUpSharedTxnRepo, setUpTrackerRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const payloadSchema = z.object({
  users: z.array(z.string()).min(2),
  name: z.string(),
});
export type CreateTrackerInput = z.infer<typeof payloadSchema>;

export default async function createTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'invalid method' });
  }

  let [parsedInput, err] = withTryCatch(() => payloadSchema.parse(req.body));
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid input' });
  }

  const session = await getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }
  const sessionUser = session.user?.email!;
  if (!parsedInput!.users.includes(sessionUser)) {
    return res
      .status(403)
      .json({ error: 'cannot create a tracker you are not a part of' });
  }

  const trackerRepo = setUpTrackerRepo();
  [err] = await withAsyncTryCatch(trackerRepo.createTracker(parsedInput!));
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong while creating a tracker' });
  }
  return res.status(202);
}
