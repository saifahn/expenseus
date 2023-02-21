import { sharedTxnItemToModel } from 'ddb/itemToModel';
import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { authOptions } from 'pages/api/auth/[...nextauth]';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
import { z, ZodError } from 'zod';

const queryStringSchema = z.object({
  from: z.coerce.number(),
  to: z.coerce.number(),
  trackerId: z.string(),
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

  const session = await getServerSession(req, res, authOptions);
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  const sharedTxnRepo = setUpSharedTxnRepo();
  var [items, err] = await withAsyncTryCatch(
    sharedTxnRepo.getTxnsByTrackerBetweenDates({
      tracker: parsed!.trackerId,
      from: parsed!.from,
      to: parsed!.to,
    }),
  );
  if (err) {
    return res.status(500).json({
      error:
        'something went wrong while getting shared transactions between dates',
    });
  }
  const txns = items?.map(sharedTxnItemToModel);
  return res.status(200).json(txns);
}
