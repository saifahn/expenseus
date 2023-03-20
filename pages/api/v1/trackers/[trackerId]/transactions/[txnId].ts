import { SubcategoryKeys } from 'data/categories';
import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { authOptions } from 'pages/api/auth/[...nextauth]';
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
  split: z.record(z.string(), z.number()).optional(),
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

  const session = getServerSession(req, res, authOptions);
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  if (req.method === 'PUT') {
    var [jsonParsed, err] = withTryCatch(() => JSON.parse(req.body));
    if (err) {
      return res.status(400).json({ error: 'error parsing payload' });
    }
    var [parsed, err] = withTryCatch(() =>
      updateSharedTxnPayloadSchema.parse({
        ...jsonParsed,
        tracker: req.query.trackerId,
        id: req.query.txnId,
      }),
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
    return res.status(202).json({});
  }

  if (req.method === 'DELETE') {
    var [jsonParsed, err] = withTryCatch(() => JSON.parse(req.body));
    if (err) {
      return res.status(400).json({ error: 'error parsing payload' });
    }
    var [delParsed, err] = withTryCatch(() =>
      deleteSharedTxnPayloadSchema.parse({
        ...jsonParsed,
        tracker: req.query.trackerId,
        txnId: req.query.txnId,
      }),
    );
    if (err instanceof ZodError) {
      return res.status(400).json({ error: 'invalid input' });
    }
    const sharedTxnRepo = setUpSharedTxnRepo();
    [, err] = await withAsyncTryCatch(
      sharedTxnRepo.deleteSharedTxn(delParsed!),
    );
    if (err) {
      return res
        .status(500)
        .json({ error: 'something went wrong while deleting shared txn' });
    }
    return res.status(202).json({});
  }
}
