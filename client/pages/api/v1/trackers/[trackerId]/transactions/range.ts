import { NextApiRequest, NextApiResponse } from 'next';

export default async function getTxnsByTrackerBetweenDatesHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).send({ error: 'invalid method' });
  }
}
