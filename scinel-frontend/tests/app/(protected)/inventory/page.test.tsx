import { render, screen } from "@testing-library/react";
import InventoryPage from "@/app/(protected)/inventory/page";
import { describe, it, expect } from "@jest/globals";

jest.mock("@/context/AuthContext", () => ({ useAuth: () => ({ user: { role: "admin" }, isAuthenticated: true }) }));

describe("InventoryPage", () => {
  it("should render", () => {
    render(<InventoryPage />);
    expect(document.body).toBeInTheDocument();
  });
});
