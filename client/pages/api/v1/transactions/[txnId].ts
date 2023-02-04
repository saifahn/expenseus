import { SubcategoryKeys } from 'data/categories';
import { txnItemToTxn } from 'ddb/itemToModel';
import { setUpDdb } from 'ddb/schema';
import { makeTxnRepository } from 'ddb/txns';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { z, ZodError } from 'zod';

const updateTxnPayloadSchema = z.object({
  id: z.string().min(1),
  userId: z.string().min(1),
  location: z.string().min(1),
  amount: z.number().min(1),
  date: z.number().min(1),
  category: SubcategoryKeys,
  details: z.string(),
});
type UpdateTxnPayload = z.infer<typeof updateTxnPayloadSchema>;

export default async function byTxnIdHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const session = await getServerSession();
  if (!session?.user) {
    res.status(401).json({ error: 'no valid session found' });
    return;
  }
  const { txnId } = req.query;
  // TODO: get ddb name from env
  const ddb = setUpDdb('test-ddb');
  const txnRepo = makeTxnRepository(ddb);

  if (req.method === 'GET') {
    const txnItem = await txnRepo.getTxn({
      userId: session.user.email!,
      txnId: txnId as string,
    });
    const item = txnItemToTxn(txnItem);
    res.status(200).json(item);
    return;
  }

  if (req.method === 'PUT') {
    let parsed: UpdateTxnPayload;
    try {
      parsed = updateTxnPayloadSchema.parse(req.body);
    } catch (err) {
      if (err instanceof ZodError) {
        res
          .status(400)
          .json({ error: 'incorrect schema for updating a transaction' });
        return;
      }
    }
    await txnRepo.updateTxn(parsed);
    res.status(202);
    return;
  }
}
