import { NextApiRequest, NextApiResponse } from 'next';

export default async function getUnsettledTxnsByTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }
}
