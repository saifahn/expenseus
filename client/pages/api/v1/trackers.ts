import { NextApiRequest, NextApiResponse } from 'next';

export default async function createTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'invalid method' });
  }
}
