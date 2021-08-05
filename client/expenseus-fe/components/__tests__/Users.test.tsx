import Users from "components/Users";
import { render, screen } from "tests/test-utils";

describe("Users component", () => {
  it("should render a Users heading", () => {
    render(<Users />);

    expect(screen.getByText("Users")).toBeInTheDocument();
  });
});
