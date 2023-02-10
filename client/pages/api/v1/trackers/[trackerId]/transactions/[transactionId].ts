import { SubcategoryKeys } from 'data/categories';
import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const updateSharedTxnPayloadSchema = z.object({
  id: z.string(),
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
export type UpdateSharedTxnPayload = z.infer<
  typeof updateSharedTxnPayloadSchema
>;

const deleteSharedTxnPayloadSchema = z.object({
  tracker: z.string(),
  txnId: z.string(),
  participants: z.array(z.string()).min(2),
});

export default async function bySharedTxnIdHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (!['PUT', 'DELETE'].includes(req.method!)) {
    return res.status(405).json({ error: 'invalid method' });
  }

  const session = getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  if (req.method === 'PUT') {
    const txn = {
      ...req.body,
      tracker: req.query.trackerId,
      id: req.query.transactionId,
    };
    let [parsed, err] = withTryCatch(() =>
      updateSharedTxnPayloadSchema.parse(txn),
    );
    if (err instanceof ZodError) {
      return res.status(400).json({ error: 'invalid input' });
    }
    const sharedTxnRepo = setUpSharedTxnRepo();
    [, err] = await withAsyncTryCatch(sharedTxnRepo.updateSharedTxn(parsed!));
    if (err) {
      return res
        .status(500)
        .json({ error: 'something went wrong while updating shared txn' });
    }
    return res.status(202);
  }

  if (req.method === 'DELETE') {
    const input = {
      ...req.body,
      tracker: req.query.trackerId,
      txnId: req.query.transactionId,
    };
    let [parsed, err] = withTryCatch(() =>
      deleteSharedTxnPayloadSchema.parse(input),
    );
    if (err instanceof ZodError) {
      return res.status(400).json({ error: 'invalid input' });
    }
    const sharedTxnRepo = setUpSharedTxnRepo();
    [, err] = await withAsyncTryCatch(sharedTxnRepo.deleteSharedTxn(parsed!));
    if (err) {
      return res
        .status(500)
        .json({ error: 'something went wrong while deleting shared txn' });
    }
    return res.status(202);
  }
}
