import { render, screen } from "tests/test-utils";
import Layout from "components/Layout";

describe("Layout component", () => {
  it("renders links to Home ", () => {
    render(<Layout>content</Layout>);

    expect(screen.getByRole("link", { name: /Home/ })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Users/ })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Expenses/ })).toBeInTheDocument();
  });
});
