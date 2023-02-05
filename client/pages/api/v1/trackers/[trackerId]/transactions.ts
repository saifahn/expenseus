import { SubcategoryKeys } from 'data/categories';
import { sharedTxnItemToModel } from 'ddb/itemToModel';
import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const createSharedTxnPayloadSchema = z.object({
  date: z.number(),
  amount: z.number(),
  location: z.string(),
  category: SubcategoryKeys,
  tracker: z.string(),
  participants: z.array(z.string()).min(2),
  payer: z.string(),
  details: z.string(),
  unsettled: z.boolean().optional(),
});
export type CreateSharedTxnPayload = z.infer<
  typeof createSharedTxnPayloadSchema
>;

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

  if (req.method === 'GET') {
    const [items, err] = await withAsyncTryCatch(
      sharedTxnRepo.getTxnsByTracker(req.query.trackerId as string),
    );
    if (err) {
      return res.status(500).json({
        error: 'something went wrong while getting transactions from tracker',
      });
    }
    const txns = items?.map(sharedTxnItemToModel);
    return res.status(200).json(txns);
  }

  if (req.method === 'POST') {
    let [parsed, err] = withTryCatch(() =>
      createSharedTxnPayloadSchema.parse(req.body),
    );
    if (err instanceof ZodError) {
      return res.status(400).json({ error: 'invalid input' });
    }

    const sessionUser = session.user?.email!;
    if (!parsed?.participants.includes(sessionUser)) {
      return res.status(403).json({
        error: 'you cannot create a shared txn without being a participant',
      });
    }

    [, err] = await withAsyncTryCatch(sharedTxnRepo.createSharedTxn(parsed));
    if (err) {
      return res
        .status(500)
        .json({ error: 'something went wrong while creating shared txn' });
    }
    return res.status(202);
  }
}
