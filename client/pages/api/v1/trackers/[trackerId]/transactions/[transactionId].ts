import { NextApiRequest, NextApiResponse } from 'next';

export default async function bySharedTxnIdHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (!['PUT', 'DELETE'].includes(req.method!)) {
    return res.status(405).json({ error: 'invalid method' });
  }
}
