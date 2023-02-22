import { sharedTxnItemToModel } from 'ddb/itemToModel';
import { setUpSharedTxnRepo } from 'ddb/setUpRepos';
import { SharedTxn } from 'ddb/sharedTxns';
import { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { authOptions } from 'pages/api/auth/[...nextauth]';
import { withAsyncTryCatch } from 'utils/withTryCatch';

/**
 * Calculates how much money is owed to the current user given a list of
 * shared transactions.
 */
function calculateDebts(currentUser: string, txns: SharedTxn[]) {
  if (txns.length === 0) return;

  // all of the txns come from the same tracker so should have the same participants
  // so can be set from the first txn
  const otherUser = txns[0].participants.find((user) => user !== currentUser)!;
  const defaultSplit = 0.5;
  let amountOwed = 0;

  for (const txn of txns) {
    let split;
    // calculate from the perspective that the logged in user is the one who has paid
    const currentUserIsPayer = txn.payer === currentUser;
    if (currentUserIsPayer) {
      // the split represents the proportion each participant will pay for a purchase
      // so when calculating the debt, it is the inverse proportion (i.e. the
      // other person's proportion) that is used
      split = txn.split?.[otherUser] ?? defaultSplit;
      amountOwed += txn.amount * split;
    } else {
      split = txn.split?.[currentUser] ?? defaultSplit;
      amountOwed -= txn.amount * split;
    }
  }

  return {
    transactions: txns,
    debtor: otherUser,
    debtee: currentUser,
    amountOwed,
  };
}

export default async function getUnsettledTxnsByTrackerHandler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method !== 'GET') {
    return res.status(405).json({ error: 'invalid method' });
  }

  const session = await getServerSession(req, res, authOptions);
  if (!session) {
    return res.status(401).json({ error: 'no valid session found' });
  }

  const sharedTxnRepo = setUpSharedTxnRepo();
  const [items, err] = await withAsyncTryCatch(
    sharedTxnRepo.getUnsettledTxnsByTracker(req.query.trackerId as string),
  );
  if (err) {
    return res
      .status(500)
      .json({ error: 'something went wrong while getting unsettled txns' });
  }
  const txns = items?.map(sharedTxnItemToModel);
  const response = calculateDebts(session.user?.email!, txns!);
  res.status(200).json(response);
}
