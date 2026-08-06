import { render, screen } from "@testing-library/react";
import NewTransactionPage from "@/app/(protected)/inventory/transaction/new/page";
import { describe, it, expect } from "@jest/globals";

jest.mock("@/context/AuthContext", () => ({ useAuth: () => ({ user: { role: "admin" }, isAuthenticated: true }) }));

describe("NewTransactionPage", () => {
  it("should render", () => {
    render(<NewTransactionPage />);
    expect(document.body).toBeInTheDocument();
  });
});
