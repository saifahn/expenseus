// this file contains the config for testing the React components
module.exports = {
  setupFilesAfterEnv: ['./react-jest-setup.ts'],
  moduleNameMapper: {
    '\\.(css)$': 'identity-obj-proxy',
  },
  moduleDirectories: ['node_modules', __dirname],
  testEnvironment: 'jsdom',
  testMatch: ['**/?(*.)+(spec|test).tsx'],
};
