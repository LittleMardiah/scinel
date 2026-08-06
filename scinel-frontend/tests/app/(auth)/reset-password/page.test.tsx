import { render, screen } from "@testing-library/react";
import ResetPasswordPage from "@/app/(auth)/reset-password/page";
import { describe, it, expect } from "@jest/globals";

describe("ResetPasswordPage", () => {
  it("should render", () => {
    render(<ResetPasswordPage />);
    expect(document.body).toBeInTheDocument();
  });
});
