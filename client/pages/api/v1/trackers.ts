import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const payloadSchema = z.object({
  users: z.array(z.string()).min(2),
  name: z.string(),
});
export type CreateTrackerInput = z.infer<typeof payloadSchema>;

export default async function createTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const [parsedInput, err] = withTryCatch(() => payloadSchema.parse(req.body));
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid input' });
  }

  const session = await getServerSession();
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }
}
