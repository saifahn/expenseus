import ExpenseList from "components/ExpenseList";
import { render, screen } from "tests/test-utils";

describe("ExpenseList component", () => {
  test("has an expense name input", () => {
    render(<ExpenseList />);

    expect(screen.getByLabelText("Name")).toBeInTheDocument();
  });

  test("has an input to upload an image", () => {
    render(<ExpenseList />);

    expect(
      screen.getByRole("button", { name: /Add picture/ })
    ).toBeInTheDocument();
  });
});
