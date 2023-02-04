import { makeTxnRepository } from 'ddb/txns';
import { NextApiRequest, NextApiResponse } from 'next';
import { RequestMethod, createMocks } from 'node-mocks-http';
import createTxnHandler, { CreateTxnPayload } from './transactions';

jest.mock('ddb/txns');
const txnsRepo = jest.mocked(makeTxnRepository);
const mockedRepoReturn: ReturnType<typeof makeTxnRepository> = {
  createTxn: jest.fn(),
  getTxn: jest.fn(),
  updateTxn: jest.fn(),
  deleteTxn: jest.fn(),
  getTxnsByUserId: jest.fn(),
  getBetweenDates: jest.fn(),
};

describe('/api/v1/transactions POST endpoint', () => {
  function mockReqRes(method: RequestMethod = 'POST') {
    const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
      method,
    });
    return { req, res };
  }

  test('should return a 405 error if the method is not POST', async () => {
    const { req, res } = mockReqRes('GET');

    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('should return a 400 error if the payload is invalid', async () => {
    const { req, res } = mockReqRes();

    req._setBody({
      invalid: 'property',
      another: 'wrong-one',
    });
    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(400);
  });

  test('should return a 200 if the payload is OK', async () => {
    const { req, res } = mockReqRes();

    const payload: CreateTxnPayload = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000 * 1000,
      category: 'unspecified.unspecified',
      details: '',
    };
    req._setBody(payload);
    txnsRepo.mockImplementationOnce(() => mockedRepoReturn);
    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(200);
  });

  test('should return a 500 if something goes wrong with ddb', async () => {
    const { req, res } = mockReqRes();

    const payload: CreateTxnPayload = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000 * 1000,
      category: 'unspecified.unspecified',
      details: '',
    };
    req._setBody(payload);
    txnsRepo.mockImplementationOnce(() => ({
      ...mockedRepoReturn,
      createTxn: jest.fn(() => {
        throw new Error();
      }),
    }));
    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(500);
  });
});
