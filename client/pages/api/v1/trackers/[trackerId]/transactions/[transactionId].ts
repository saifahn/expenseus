import { SubcategoryKeys } from 'data/categories';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withTryCatch } from 'utils/withTryCatch';
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

  let [parsed, err] = withTryCatch(() =>
    updateSharedTxnPayloadSchema.parse(req.body),
  );
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid input' });
  }
}
