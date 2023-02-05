import { NextApiRequest, NextApiResponse } from 'next';

export default async function txnsByTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (!['GET', 'POST'].includes(req.body)) {
    return res.status(405).json({ error: 'invalid method' });
  }
}
