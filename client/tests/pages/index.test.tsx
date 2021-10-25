import React from "react";
import { rest } from "msw";
import { setupServer } from "msw/node";
import Home from "pages/index";
import { render, screen, waitFor } from "tests/test-utils";
import { testSeanUser } from "tests/doubles/index";

const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

const apiBaseURL = process.env.NEXT_PUBLIC_API_BASE_URL;

describe("HomePage", () => {
  describe("when not logged in", () => {
    server.use(
      rest.get(`${apiBaseURL}/users/self`, (_, res, ctx) => {
        return res(ctx.status(401), ctx.json(""));
      })
    );
  });

  it("should show a 'Sign in with Google' button if not logged in and not show a 'Log out' button", async () => {
    render(<Home />);

    const signInButton = await waitFor(() =>
      screen.getByRole("link", { name: /Sign in with Google/ })
    );
    expect(signInButton).toBeInTheDocument();

    const logOutButton = screen.queryByRole("link", { name: /Log out/ });
    expect(logOutButton).not.toBeInTheDocument();
  });

  describe("when logged in", () => {
    beforeEach(() => {
      server.use(
        rest.get(`${apiBaseURL}/users/self`, (_, res, ctx) => {
          return res(ctx.status(202), ctx.json(testSeanUser));
        })
      );
    });

    it("should not show a 'Sign in with Google' button, but should show a 'Log out' button", async () => {
      render(<Home />);

      // wait for the title to appear
      await waitFor(() => screen.getByText("Welcome to Expenseus"));

      const signInButton = screen.queryByRole("link", {
        name: /Sign in with Google/,
      });
      expect(signInButton).not.toBeInTheDocument();

      const logOutButton = screen.getByRole("link", { name: /Log out/ });
      expect(logOutButton).toBeInTheDocument();
    });

    it("should show a welcome message with the user's username", async () => {
      render(<Home />);

      await waitFor(() => screen.getByText("Welcome to Expenseus"));

      const welcomeText = screen.getByTestId("welcome");

      expect(welcomeText).toHaveTextContent(testSeanUser.username);
    });
  });
});
