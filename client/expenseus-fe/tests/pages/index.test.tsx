import React from "react";
import Home from "pages/index";
import { render, screen } from "tests/test-utils";

describe("HomePage", () => {
  it("should render the heading", () => {
    render(<Home />);

    const heading = screen.getByText("Welcome to Expenseus");

    expect(heading).toBeInTheDocument();
  });
});
