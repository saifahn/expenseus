import { txnItemToTxn } from 'ddb/itemToModel';
import { setUpDdb } from 'ddb/schema';
import { makeTxnRepository } from 'ddb/txns';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';

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
    await txnRepo.updateTxn(req.body);
    res.status(202);
    return;
  }
}
