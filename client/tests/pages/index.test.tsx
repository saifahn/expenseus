import React from "react";
import Home from "pages/index";
import { render, screen } from "tests/test-utils";
import { setupServer } from "msw/node";
import { rest } from "msw";

const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("HomePage", () => {
  it("should render the heading", () => {
    render(<Home />);

    const heading = screen.getByText("Welcome to Expenseus");

    expect(heading).toBeInTheDocument();
  });

  it("should show a login button if not logged in", async () => {
    server.use(
      rest.get(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/self_details`,
        (_req, res, ctx) => {
          return res(ctx.status(401), ctx.json({ err: "unauthenticated" }));
        }
      )
    );

    render(<Home />);

    const button = screen.getByRole("button", { name: "Sign in with Google" });

    expect(button).toBeInTheDocument();
  });
});
