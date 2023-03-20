import { rest } from "msw";
import { setupServer } from "msw/node";
import UserList from "components/UserList";
import { testSeanUser, testTomomiUser } from "tests/doubles/index";
import { render, screen, userEvent, waitFor } from "tests/test-utils";

const testUsers = [testSeanUser, testTomomiUser];

const server = setupServer(
  rest.get(`${process.env.NEXT_PUBLIC_API_BASE_URL}/users`, (_, res, ctx) => {
    return res(ctx.json(testUsers));
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
    render(<UserList />);

    expect(screen.getByText("Users")).toBeInTheDocument();
  });

  it("should render a list of users after loading", async () => {
    render(<UserList />);

    expect(await screen.findByText(testTomomiUser.name)).toBeInTheDocument();
    expect(await screen.findByText(testSeanUser.name)).toBeInTheDocument();
  });
});
