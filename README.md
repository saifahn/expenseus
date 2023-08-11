# Expenseus

## Dependencies

- Node
- npm or yarn
- DynamoDB

## Environment variables setup

- Add a `.env.local` file in the root folder with your
  - `GOOGLE_CLIENT_ID`
  - `GOOGLE_CLIENT_SECRET`
- Depending on your DynamoDB setup, you may also set
  - `DDB_TABLE`
  - `DDB_ENDPOINT`
  - `DDB_ACCESS_KEY_ID`
  - `DDB_SECRET_ACCESS_KEY`

## Running the application

- Make sure you have a local dynamodb

```sh
# install dependencies
yarn

# create a local dynamodb table
yarn ddb:setup

# run the dev server
yarn dev
```

You can see the UI at http://127.0.0.1:3000
