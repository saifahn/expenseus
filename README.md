# Expenseus

## Dependencies

- Go
- Redis
- Node
- npm or yarn

## Running a development version

### Front end

```sh
cd client

# install dependencies
yarn

# run the dev server
yarn dev
```

You can see the UI at http://127.0.0.1:3000

### Back end

1. Create an environment variable file with the variables found in `.env.example`
   - e.g. if you are using `direnv`, you can copy the file to a `.envrc`
2. Create a Google OAuth client on [Google Cloud Platform](https://console.cloud.google.com/)
3. Fill in the details. GOOGLE_REDIRECT_URL with the default router settings should be `http://127.0.0.1:4000/api/v1/callback_google`
4. REDIS_ADDRESS is set by default to `127.0.0.1:6379`
5. After setting up the environment variables, build the webserver and run it

```sh
# from the `server` directory
go build ./cmd/webserver

# run the application
./webserver
```

## Serving the bundle from the Go application

1. Build the front end

   ```sh
   cd client

   # build a static version of the app, exported to `out` directory
   yarn build

   # create the `web` directory in the `server` directory if it doesn't exist
   mkdir ../server/web

   # move the built static app to the `web/dist` directory
   mv out ../server/web/dist
   ```

2. Build the web server and run it

   ```sh
   # from the `server` directory
   go build ./cmd/webserver

   # run the application
   ./webserver
   ```

Navigate to http://127.0.0.1:5000 to see the UI hosted in the Go application.
