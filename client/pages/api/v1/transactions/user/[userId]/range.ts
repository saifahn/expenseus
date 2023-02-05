import { NextApiRequest, NextApiResponse } from 'next';
import { withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const queryStringSchema = z.object({
  from: z.coerce.number(),
  to: z.coerce.number(),
});

export default async function getTxnsByUserIdBetweenDatesHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const [parsed, err] = withTryCatch(() => queryStringSchema.parse(req.query));
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid query' });
  }
}
