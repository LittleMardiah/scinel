import { render, screen } from "@testing-library/react";
import AnomalyDetailPage from "@/app/(protected)/anomalies/[id]/page";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

// Tidak ada mock lokal - pakai global dari setup
describe("AnomalyDetailPage", () => {
  it("should render", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <AnomalyDetailPage />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
