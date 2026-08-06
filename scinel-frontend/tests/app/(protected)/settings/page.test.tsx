import { render, screen } from "@testing-library/react";
import SettingsPage from "@/app/(protected)/settings/page";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

describe("SettingsPage", () => {
  it("should render", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <SettingsPage />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
