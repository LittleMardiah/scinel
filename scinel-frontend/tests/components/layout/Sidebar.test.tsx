import { render, screen } from "@testing-library/react";
import { Sidebar } from "@/components/layout/Sidebar";
import { AuthProvider } from "@/context/AuthContext";
import { describe, it, expect, jest } from "@jest/globals";

jest.mock("next/navigation", () => ({
  usePathname: () => "/dashboard",
}));

// Mock useAuth agar tidak error
jest.mock("@/context/AuthContext", () => ({
  useAuth: () => ({
    user: { role: "admin" },
    isAuthenticated: true,
  }),
  AuthProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

describe("Sidebar", () => {
  it("should render navigation links", () => {
    render(
      <AuthProvider>
        <Sidebar />
      </AuthProvider>
    );
    expect(screen.getByText(/Dashboard/i)).toBeInTheDocument();
    expect(screen.getByText(/Inventory/i)).toBeInTheDocument();
  });
});
