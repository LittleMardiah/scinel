import { render, screen } from "@testing-library/react";
import TransactionsPage from "@/app/(protected)/inventory/transactions/page";
import { describe, it, expect } from "@jest/globals";

jest.mock("@/context/AuthContext", () => ({ useAuth: () => ({ user: { role: "admin" }, isAuthenticated: true }) }));

describe("TransactionsPage", () => {
  it("should render", () => {
    render(<TransactionsPage />);
    expect(document.body).toBeInTheDocument();
  });
});
