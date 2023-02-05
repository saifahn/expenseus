import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const payloadSchema = z.array(
  z.object({
    id: z.string(),
    trackerId: z.string(),
    participants: z.array(z.string()),
  }),
);

export default async function settleTxnsHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const [parsed, err] = withTryCatch(() => payloadSchema.parse(req.body));
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid input' });
  }

  const session = await getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }
}
