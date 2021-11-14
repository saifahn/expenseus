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
      render(<UserList />);

      const nameInput = screen.getByRole("textbox", { name: /Name/ });
      userEvent.type(nameInput, "Test User");

      const usernameInput = screen.getByRole("textbox", { name: /Username/ });
      userEvent.type(usernameInput, "testuser");

      const submitButton = screen.getByRole("button", { name: /Create user/ });
      userEvent.click(submitButton);

      await waitFor(() =>
        expect(screen.getByRole("status")).toHaveTextContent(
          `User testuser was successfully created`
        )
      );
    });

    it("should add a different testuser and render it", async () => {
      server.use(
        rest.post(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`,
          (req, res, ctx) => {
            const { username, name, id } = req.body as {
              username: string;
              name: string;
              id: string;
            };
            const newUser = {
              username,
              name,
              id,
            };
            testUsers.push(newUser);
            return res(ctx.status(202));
          }
        )
      );
      render(<UserList />);

      const nameInput = screen.getByRole("textbox", { name: /Name/ });
      userEvent.type(nameInput, "Test Usertwo");

      const usernameInput = screen.getByRole("textbox", { name: /Username/ });
      userEvent.type(usernameInput, "testuser2");

      const submitButton = screen.getByRole("button", { name: /Create user/ });
      userEvent.click(submitButton);

      await waitFor(() =>
        expect(screen.getByRole("status")).toHaveTextContent(
          `User testuser2 was successfully created`
        )
      );

      await waitFor(() =>
        expect(screen.getByText("Test Usertwo")).toBeInTheDocument()
      );
    });
  });
});
