module.exports = {
  setupFilesAfterEnv: ['./api-jest-setup.ts'],
  moduleDirectories: ['node_modules', __dirname],
  testMatch: ['**/?(*.)+(spec|test).ts'],
};
