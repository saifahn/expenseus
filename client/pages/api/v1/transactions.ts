import { NextApiRequest, NextApiResponse } from 'next';
import { z, ZodError } from 'zod';

const createTxnPayloadSchema = z.object({
  userId: z.string(),
  location: z.string(),
  amount: z.number(),
  date: z.number(),
  category: z.string(),
  details: z.string(),
});

export default async function createTxnHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    res.status(405).json({ error: 'invalid method' });
    return;
  }

  try {
    createTxnPayloadSchema.parse(req.body);
  } catch (err) {
    if (err instanceof ZodError) {
      res.status(400).json({ error: 'invalid payload' });
      return;
    }
  }
}
