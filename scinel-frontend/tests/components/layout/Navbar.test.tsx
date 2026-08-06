import { render, screen } from "@testing-library/react";
import { Navbar } from "@/components/layout/Navbar";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

describe("Navbar", () => {
  it("should render logo", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <Navbar />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(screen.getByText(/ScIntel/i)).toBeInTheDocument();
  });

  it("should render theme toggle", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <Navbar />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(screen.getByTestId("theme-toggle")).toBeInTheDocument();
  });
});
