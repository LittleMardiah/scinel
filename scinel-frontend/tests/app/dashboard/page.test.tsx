import { render, screen } from "@testing-library/react";
import DashboardPage from "@/app/(protected)/dashboard/page";
import { describe, it, expect, jest } from "@jest/globals";

// Mock useDashboardData dengan data health yang lengkap
jest.mock("@/hooks/useDashboardData", () => ({
  useDashboardData: () => ({
    data: {
      total_sku: 100,
      total_warehouse: 3,
      low_stock_count: 5,
      anomaly_count: 2,
      health: { status: "online", message: "All systems operational" },
      recent_transactions: [],
      forecast: [],
    },
    isLoading: false,
    error: null,
  }),
}));

// Mock useAuth agar user terautentikasi
jest.mock("@/context/AuthContext", () => ({
  useAuth: () => ({
    user: { id: "1", name: "Test User" },
    isAuthenticated: true,
  }),
}));

describe("DashboardPage", () => {
  it("should render dashboard", () => {
    render(<DashboardPage />);
    expect(screen.getByText(/Dashboard/i)).toBeInTheDocument();
  });
});
