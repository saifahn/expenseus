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

describe("HomePage", () => {
  it("should show a 'Sign in with Google' button if not logged in", async () => {
    server.use(
      rest.get(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/users/self`,
        (_, res, ctx) => {
          return res(ctx.status(401), ctx.json(""));
        }
      )
    );

    render(<Home />);

    const button = await waitFor(() =>
      screen.getByRole("button", { name: "Sign in with Google" })
    );

    expect(button).toBeInTheDocument();
  });

  describe("when logged in", () => {
    beforeEach(() => {
      server.use(
        rest.get(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/users/self`,
          (_, res, ctx) => {
            return res(ctx.status(202), ctx.json(testSeanUser));
          }
        )
      );
    });

    it("should not show a 'Sign in with Google' button", async () => {
      render(<Home />);

      // wait for the title to appear
      await waitFor(() => screen.getByText("Welcome to Expenseus"));

      const button = screen.queryByRole("button", {
        name: "Sign in with Google",
      });

      expect(button).not.toBeInTheDocument();
    });

    it("should show a welcome message with the user's username", async () => {
      render(<Home />);

      await waitFor(() => screen.getByText("Welcome to Expenseus"));

      const welcomeText = screen.getByTestId("welcome");

      expect(welcomeText).toHaveTextContent(testSeanUser.username);
    });
  });
});
