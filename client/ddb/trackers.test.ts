import { setUpDdb, createTableIfNotExists, deleteTable } from 'ddb/schema';
import * as ulid from 'ulid';
import { createTracker, getTracker, getTrackersByUser } from './trackers';

jest.mock('ulid');
let mockedUlid = jest.mocked(ulid);

const trackersTestTable = 'trackers-test-table';
const d = setUpDdb(trackersTestTable);

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
    await createTracker(d, {
      users: ['user-01', 'user-02'],
      name: 'The Test Tracker',
    });

    // get the tracker based on the above ulid
    const tracker = await getTracker(d, testUlid);
    expect(tracker).toBeTruthy();

    const notTracker = await getTracker(d, 'non-existent');
    expect(notTracker).toBeFalsy();

    const user01Trackers = await getTrackersByUser(d, 'user-01');
    expect(user01Trackers).toHaveLength(1);

    const user02Trackers = await getTrackersByUser(d, 'user-02');
    expect(user02Trackers).toHaveLength(1);

    const notUserTrackers = await getTrackersByUser(d, 'not-user');
    expect(notUserTrackers).toHaveLength(0);

    mockedUlid.ulid.mockReset();
  });
});
