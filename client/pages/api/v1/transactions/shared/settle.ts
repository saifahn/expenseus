import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const payloadSchema = z.array(
  z.object({
    id: z.string(),
    trackerId: z.string(),
    participants: z.array(z.string()),
  }),
);

export default async function settleTxnsHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'invalid method' });
  }

  let [parsed, err] = withTryCatch(() => payloadSchema.parse(req.body));
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid input' });
  }

  const session = await getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }
  const sessionUser = session.user?.email!;

  for (const txn of parsed!) {
    if (!txn.participants.includes(sessionUser)) {
      return res.status(403).json({
        error: 'you cannot settle a transaction you do not belong to',
      });
    }
  }

  const sharedTxnRepo = setUpSharedTxnRepo();
  [err] = await withAsyncTryCatch(sharedTxnRepo.settleTxns(parsed!));
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong while settling transactions' });
  }
  return res.status(202);
}
