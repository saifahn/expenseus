import { NextApiRequest, NextApiResponse } from 'next';
import { RequestMethod, createMocks } from 'node-mocks-http';
import createTxnHandler from './transactions';

describe('/api/v1/transactions POST endpoint', () => {
  function mockReqRes(method: RequestMethod = 'POST') {
    const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
      method,
    });
    return { req, res };
  }

  it('should return a 405 error if the method is not POST', async () => {
    const { req, res } = mockReqRes('GET');

    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  it('should return a 400 error if the payload is invalid', async () => {
    const { req, res } = mockReqRes();

    req._setBody({
      invalid: 'property',
      another: 'wrong-one',
    });
    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(400);
  });
});
