import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import * as ulid from 'ulid';
import { makeTrackerRepository } from './trackers';

jest.mock('ulid');
let mockedUlid = jest.mocked(ulid);

const trackersTestTable = 'trackers-test-table';
const d = setUpDdb(trackersTestTable);
const { createTracker, getTracker, getTrackersByUser } =
  makeTrackerRepository(d);

describe('Trackers', () => {
  beforeEach(async () => {
    await createTableIfNotExists(trackersTestTable);
  });

  afterEach(async () => {
    await deleteTable(trackersTestTable);
  });

  test('a tracker can be created and retrieved successfully', async () => {
    // define the ulid that will be generated
    const testUlid = '1234';
    mockedUlid.ulid.mockImplementationOnce(() => testUlid);
    await createTracker({
      users: ['user-01', 'user-02'],
      name: 'The Test Tracker',
    });

    // get the tracker based on the above ulid
    const tracker = await getTracker(testUlid);
    expect(tracker).toBeTruthy();

    const notTracker = await getTracker('non-existent');
    expect(notTracker).toBeFalsy();

    const user01Trackers = await getTrackersByUser('user-01');
    expect(user01Trackers).toHaveLength(1);

    const user02Trackers = await getTrackersByUser('user-02');
    expect(user02Trackers).toHaveLength(1);

    const notUserTrackers = await getTrackersByUser('not-user');
    expect(notUserTrackers).toHaveLength(0);

    mockedUlid.ulid.mockReset();
  });
});
