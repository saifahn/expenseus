import { makeTxnRepository } from 'ddb/txns';
import { mockReqRes } from 'tests/api/common';
import createTxnHandler, { CreateTxnPayload } from './transactions';

jest.mock('ddb/txns');
const txnsRepo = jest.mocked(makeTxnRepository);
export const txnRepoFnsMock: ReturnType<typeof makeTxnRepository> = {
  createTxn: jest.fn(),
  getTxn: jest.fn(),
  updateTxn: jest.fn(),
  deleteTxn: jest.fn(),
  getTxnsByUserId: jest.fn(),
  getBetweenDates: jest.fn(),
};

describe('/api/v1/transactions POST endpoint', () => {
  test('should return a 405 error if the method is not POST', async () => {
    const { req, res } = mockReqRes('GET');

    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(405);
  });

  test('should return a 400 error if the payload is invalid', async () => {
    const { req, res } = mockReqRes('POST');

    req._setBody({
      invalid: 'property',
      another: 'wrong-one',
    });
    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(400);
  });

  test('should return a 202 if the payload is OK', async () => {
    const { req, res } = mockReqRes('POST');

    const payload: CreateTxnPayload = {
      userId: 'test-user',
      location: 'test-location',
      amount: 12345,
      date: 1000 * 1000,
      category: 'unspecified.unspecified',
      details: '',
    };
    req._setBody(payload);
    txnsRepo.mockImplementationOnce(() => txnRepoFnsMock);
    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(202);
  });

  test('should return a 500 if something goes wrong with ddb', async () => {
    const { req, res } = mockReqRes('POST');

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
      ...txnRepoFnsMock,
      createTxn: jest.fn(async () => {
        throw new Error();
      }),
    }));
    await createTxnHandler(req, res);

    expect(res.statusCode).toBe(500);
  });
});