import { render, screen } from "@testing-library/react";
import AnomaliesPage from "@/app/(protected)/anomalies/page";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

describe("AnomaliesPage", () => {
  it("should render", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <AnomaliesPage />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
