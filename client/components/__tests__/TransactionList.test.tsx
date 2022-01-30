import TransactionList from "components/TransactionList";
import { render, screen } from "tests/test-utils";

describe("TransactionList component", () => {
  test("has an transaction name input", () => {
    render(<TransactionList />);

    expect(screen.getByLabelText("Name")).toBeInTheDocument();
  });

  test("has an input to upload an image", () => {
    render(<TransactionList />);

    expect(
      screen.getByRole("button", { name: /Add picture/ })
    ).toBeInTheDocument();
  });
});
