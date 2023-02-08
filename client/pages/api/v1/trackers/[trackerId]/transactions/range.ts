import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const queryStringSchema = z.object({
  from: z.coerce.number(),
  to: z.coerce.number(),
});

export default async function getTxnsByTrackerBetweenDatesHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).send({ error: 'invalid method' });
  }

  var [parsed, err] = withTryCatch(() => queryStringSchema.parse(req.query));
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid query' });
  }

  const session = await getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }
}
