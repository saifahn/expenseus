import dotenv from 'dotenv';

dotenv.config({ path: '.env.test' });
jest.mock('next-auth');
