import { render, screen } from "@testing-library/react";
import { InventoryFilters } from "@/components/inventory/InventoryFilters";
import { describe, it, expect, jest } from "@jest/globals";

describe("InventoryFilters", () => {
  it("should render filters", () => {
    render(<InventoryFilters onFilterChange={jest.fn()} />);
    
    // Cari label secara spesifik (hanya label, bukan option)
    const gudangLabel = screen.getByText(/Gudang/i, { selector: "label" });
    expect(gudangLabel).toBeInTheDocument();
    
    const statusLabel = screen.getByText(/Status/i, { selector: "label" });
    expect(statusLabel).toBeInTheDocument();
    
    // Verifikasi ada 2 combobox (select)
    expect(screen.getAllByRole("combobox").length).toBe(2);
  });
});
