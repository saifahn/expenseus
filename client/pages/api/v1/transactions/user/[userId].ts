import { NextApiRequest, NextApiResponse } from 'next';

export default async function txnByUserIdHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }
}
