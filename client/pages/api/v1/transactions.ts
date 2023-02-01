import { NextApiRequest, NextApiResponse } from 'next';
import { z, ZodError } from 'zod';

const createTxnPayloadSchema = z.object({
  userId: z.string().min(1),
  location: z.string().min(1),
  amount: z.number().min(1),
  date: z.number().min(1),
  category: z.string().min(1),
  details: z.string(),
});

export type CreateTxnPayload = z.infer<typeof createTxnPayloadSchema>;

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
