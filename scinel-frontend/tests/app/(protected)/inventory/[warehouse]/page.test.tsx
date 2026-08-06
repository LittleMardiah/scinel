import { render, screen } from "@testing-library/react";
import WarehouseInventoryPage from "@/app/(protected)/inventory/[warehouse]/page";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

// Tidak ada mock lokal - pakai global dari setup
describe("WarehouseInventoryPage", () => {
  it("should render", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <WarehouseInventoryPage />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
