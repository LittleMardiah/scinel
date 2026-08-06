import { render, screen } from "@testing-library/react";
import LoginPage from "@/app/(auth)/login/page";
import { AuthProvider } from "@/context/AuthContext";
import { describe, it, expect } from "@jest/globals";

describe("LoginPage", () => {
  it("should render", () => {
    render(
      <AuthProvider>
        <LoginPage />
      </AuthProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
