import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { LoginForm } from "@/components/auth/LoginForm";
import { describe, it, expect, jest } from "@jest/globals";

describe("LoginForm", () => {
  it("should render email and password fields", () => {
    render(<LoginForm onSubmit={jest.fn(async () => {})} />);
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/kata sandi/i)).toBeInTheDocument();
  });

  it("should validate email format", async () => {
    const handleSubmit = jest.fn(async () => {});
    render(<LoginForm onSubmit={handleSubmit} />);
    const emailInput = screen.getByLabelText(/email/i);
    const submitBtn = screen.getByRole("button", { name: /masuk/i });
    const form = submitBtn.closest("form") as HTMLFormElement;

    fireEvent.change(emailInput, { target: { value: "invalid-email" } });
    fireEvent.submit(form);

    // Ada dua elemen error (dari Input dan dari LoginForm), jadi gunakan getAllByText
    const errorMessages = await screen.findAllByText(/format email tidak valid/i);
    expect(errorMessages.length).toBeGreaterThan(0);
    expect(handleSubmit).not.toHaveBeenCalled();
  });

  it("should submit with valid credentials", async () => {
    const handleSubmit = jest.fn<(email: string, password: string) => Promise<void>>(async () => {});
    render(<LoginForm onSubmit={handleSubmit} />);
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: "test@example.com" } });
    fireEvent.change(screen.getByLabelText(/kata sandi/i), { target: { value: "password123" } });
    fireEvent.click(screen.getByRole("button", { name: /masuk/i }));
    expect(handleSubmit).toHaveBeenCalledWith("test@example.com", "password123");
  });

  it("should show error message from prop", () => {
    render(<LoginForm onSubmit={jest.fn(async () => {})} error="Login failed" />);
    expect(screen.getByText("Login failed")).toBeInTheDocument();
  });
});
