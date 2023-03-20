import { txnItemToTxn } from 'ddb/itemToModel';
import { setUpTxnRepo } from 'ddb/setUpRepos';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { authOptions } from 'pages/api/auth/[...nextauth]';
import { withAsyncTryCatch, withTryCatch } from 'utils/withTryCatch';
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

  var [parsed, err] = withTryCatch(() => queryStringSchema.parse(req.query));
  if (err instanceof ZodError) {
    return res.status(400).json({ error: 'invalid query' });
  }

  const session = await getServerSession(req, res, authOptions);
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  const txnRepo = setUpTxnRepo();
  var [items, err] = await withAsyncTryCatch(
    txnRepo.getBetweenDates({
      userId: session.user!.email!,
      from: parsed!.from,
      to: parsed!.to,
    }),
  );
  if (err) {
    return res.status(500).json({
      error:
        'something went wrong while getting transactions between dates for a user',
    });
  }
  const txns = items?.map(txnItemToTxn);
  return res.status(200).json(txns);
}
