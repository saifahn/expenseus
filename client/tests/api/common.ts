import { NextApiRequest, NextApiResponse } from 'next';
import { RequestMethod, createMocks } from 'node-mocks-http';

export function mockReqRes(method: RequestMethod) {
  const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
    method,
  });
  return { req, res };
}
