import { render, screen } from "@testing-library/react";
import { RecentTransactions } from "@/components/dashboard/RecentTransactions";
import { describe, it, expect } from "@jest/globals";

const mockTransactions = [
  { id: "1", sku: "SKU-001", type: "in" as const, quantity: 10, created_at: "2026-07-29" },
  { id: "2", sku: "SKU-002", type: "out" as const, quantity: 5, created_at: "2026-07-28" },
];

describe("RecentTransactions", () => {
  it("should render transactions", () => {
    render(<RecentTransactions transactions={mockTransactions} />);
    expect(screen.getByText("SKU-001")).toBeInTheDocument();
    expect(screen.getByText("SKU-002")).toBeInTheDocument();
  });

  it("should show empty state", () => {
    render(<RecentTransactions transactions={[]} />);
    expect(screen.getByText(/Belum ada transaksi/i)).toBeInTheDocument();
  });
});
