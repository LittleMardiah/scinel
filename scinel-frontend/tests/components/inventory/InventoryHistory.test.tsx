import { render, screen } from "@testing-library/react";
import { InventoryHistory } from "@/components/inventory/InventoryHistory";
import { describe, it, expect } from "@jest/globals";

describe("InventoryHistory", () => {
  it("should render", () => {
    render(<InventoryHistory transactions={[]} />);
    expect(document.body).toBeInTheDocument();
  });
});
