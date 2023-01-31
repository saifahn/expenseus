import { NextApiRequest, NextApiResponse } from 'next';
import { createMocks, RequestMethod } from 'node-mocks-http';
import usersHandler from '../users';

describe('/api/v1/users API endpoint', () => {
  function mockReqRes(method: RequestMethod = 'GET') {
    const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
      method,
    });
    return { req, res };
  }

  it('should return a successful response', async () => {
    const { req, res } = mockReqRes();
    await usersHandler(req, res);

    expect(res.statusCode).toBe(200);
    expect(res.getHeaders()).toEqual({ 'content-type': 'application/json' });
  });
});
