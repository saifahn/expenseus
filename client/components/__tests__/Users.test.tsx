import { render, screen, userEvent, waitFor } from "tests/test-utils";
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

const testUsers = [testSeanUser, testTomomiUser];

const server = setupServer(
  rest.get(`${process.env.NEXT_PUBLIC_API_BASE_URL}/users`, (req, res, ctx) => {
    return res(ctx.json({ users: testUsers }));
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

  describe("adding new users", () => {
    it("should add a new testuser", async () => {
      server.use(
        rest.post(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`,
          (req, res, ctx) => {
            return res(ctx.status(202));
          }
        )
      );
      render(<Users />);

      const nameInput = screen.getByRole("textbox", { name: /Name/ });
      userEvent.type(nameInput, "Test User");

      const usernameInput = screen.getByRole("textbox", { name: /Username/ });
      userEvent.type(usernameInput, "testuser");

      const submitButton = screen.getByRole("button", { name: /Create user/ });
      userEvent.click(submitButton);

      await waitFor(() =>
        expect(screen.getByRole("status")).toHaveTextContent(
          `User testuser successfully created`
        )
      );
    });

    // NOTE: I think this isn't such a good pattern as I'm mocking the behaviour of the backend too?
    it("should add a different testuser and render it", async () => {
      server.use(
        rest.post(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`,
          (req, res, ctx) => {
            const { username, name } = req.body as {
              username: string;
              name: string;
            };
            const newUser = {
              username,
              name,
              id: "test_id",
            };
            testUsers.push(newUser);
            return res(ctx.status(202));
          }
        )
      );
      render(<Users />);

      const nameInput = screen.getByRole("textbox", { name: /Name/ });
      userEvent.type(nameInput, "Test Usertwo");

      const usernameInput = screen.getByRole("textbox", { name: /Username/ });
      userEvent.type(usernameInput, "testuser2");

      const submitButton = screen.getByRole("button", { name: /Create user/ });
      userEvent.click(submitButton);

      await waitFor(() =>
        expect(screen.getByRole("status")).toHaveTextContent(
          `User testuser2 successfully created`
        )
      );

      expect(await screen.findByText("Test Usertwo")).toBeInTheDocument();
    });
  });
});
