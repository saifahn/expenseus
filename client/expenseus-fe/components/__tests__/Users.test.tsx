import { render, screen } from "tests/test-utils";
import { rest } from "msw";
import { setupServer } from "msw/node";
import Users, { User } from "components/Users";

const testSeanUser: User = {
  username: "saifahn",
  name: "Sean Li",
  id: "sean_id",
};

const testTomomiUser: User = {
  username: "tomochi",
  name: "Tomomi Kinoshita",
  id: "tomomi_id",
};

const server = setupServer(
  rest.get(`${process.env.API_BASE_URL}/users`, (req, res, ctx) => {
    return res(ctx.json({ users: [testSeanUser, testTomomiUser] }));
  })
);

// Enable API mocking before tests
beforeAll(() => server.listen());

// Reset any runtime request handlers we may add during the tests
afterEach(() => server.resetHandlers());

// Disable API mocking after the tests are done
afterAll(() => server.close());

describe("Users component", () => {
  it("should render a Users heading", () => {
    render(<Users />);

    expect(screen.getByText("Users")).toBeInTheDocument();
  });

  it("should render a list of users after loading", async () => {
    render(<Users />);

    expect(await screen.findByText(testTomomiUser.name)).toBeInTheDocument();
    expect(await screen.findByText(testSeanUser.name)).toBeInTheDocument();
  });
});
