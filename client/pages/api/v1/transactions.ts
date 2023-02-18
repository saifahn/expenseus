/**
 * This is the file for the CreateTxnHandler
 */
import { SubcategoryKeys } from 'data/categories';
import { setUpTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const createTxnPayloadSchema = z.object({
  userId: z.string().min(1),
  location: z.string().min(1),
  amount: z.number().min(1),
  date: z.number().min(1),
  category: SubcategoryKeys,
  details: z.string(),
});

export type CreateTxnPayload = z.infer<typeof createTxnPayloadSchema>;

export default async function createTxnHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'invalid method' });
  }

  let [parsed, err] = withTryCatch(() =>
    createTxnPayloadSchema.parse(JSON.parse(req.body)),
  );
  if (err) {
    return res.status(400).json({ error: 'invalid payload' });
  }

  const txnRepo = setUpTxnRepo();

  [, err] = await withAsyncTryCatch(txnRepo.createTxn(parsed!));
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong in creating the transaction' });
  }
  return res.status(202).json({});
}
