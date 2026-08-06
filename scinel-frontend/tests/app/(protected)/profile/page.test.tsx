import { render, screen } from "@testing-library/react";
import ProfilePage from "@/app/(protected)/profile/page";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

describe("ProfilePage", () => {
  it("should render", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <ProfilePage />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
