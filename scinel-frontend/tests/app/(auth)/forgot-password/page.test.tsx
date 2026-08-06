import { render, screen } from "@testing-library/react";
import ForgotPasswordPage from "@/app/(auth)/forgot-password/page";
import { describe, it, expect } from "@jest/globals";

describe("ForgotPasswordPage", () => {
  it("should render", () => {
    render(<ForgotPasswordPage />);
    expect(document.body).toBeInTheDocument();
  });
});
