import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';

export default async function getUnsettledTxnsByTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const session = getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }
}
