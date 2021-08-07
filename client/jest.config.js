module.exports = {
  setupFilesAfterEnv: ["./jest-setup.ts"],
  moduleNameMapper: {
    "\\.(css)$": "identity-obj-proxy",
  },
  moduleDirectories: [
    "node_modules",
    __dirname
  ],
  testEnvironment: "jsdom"
}