import { SubcategoryKeys } from 'data/categories';
import { txnItemToTxn } from 'ddb/itemToModel';
import { setUpTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { authOptions } from 'pages/api/auth/[...nextauth]';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
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
export type UpdateTxnPayload = z.infer<typeof updateTxnPayloadSchema>;

export default async function byTxnIdHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (!['GET', 'PUT', 'DELETE'].includes(req.method ?? '')) {
    return res.status(405).json({ error: 'method not allowed' });
  }

  const session = await getServerSession(req, res, authOptions);
  if (!session?.user) {
    return res.status(401).json({ error: 'no valid session found' });
  }
  const txnId = req.query.txnId as string;
  const sessionUser = session.user.email!;
  const txnRepo = setUpTxnRepo();

  // get transaction
  if (req.method === 'GET') {
    const txnItem = await txnRepo.getTxn({
      userId: sessionUser,
      txnId,
    });
    const item = txnItemToTxn(txnItem);
    return res.status(200).json(item);
  }

  // update transaction
  if (req.method === 'PUT') {
    var [jsonParsed, err] = withTryCatch(() => JSON.parse(req.body));
    if (err) {
      return res.status(400).json({ error: 'error parsing payload' });
    }
    var [parsed, err] = withTryCatch(() =>
      updateTxnPayloadSchema.parse({
        id: txnId,
        userId: sessionUser,
        ...jsonParsed,
      }),
    );
    if (err instanceof ZodError) {
      return res
        .status(400)
        .json({ error: 'incorrect schema for updating a transaction' });
    }

    if (parsed?.userId !== sessionUser) {
      return res.status(403).json({
        error: "you don't have permissions to update this transaction ",
      });
    }

    [, err] = await withAsyncTryCatch(txnRepo.updateTxn(parsed!));
    if (err) {
      return res
        .status(500)
        .json({ error: 'something went wrong while updating the transaction' });
    }

    return res.status(202).json({});
  }

  // delete transaction
  if (req.method === 'DELETE') {
    const [, err] = await withAsyncTryCatch(
      txnRepo.deleteTxn({
        txnId,
        userId: sessionUser,
      }),
    );
    if (err) {
      return res
        .status(500)
        .json({ error: 'something went wrong while deleting the transaction' });
    }

    return res.status(202).json({});
  }
}
