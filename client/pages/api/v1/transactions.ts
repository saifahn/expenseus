import { NextApiRequest, NextApiResponse } from 'next';

export default async function createTxnHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    res.status(405).json({ error: 'invalid method' });
  }
}
