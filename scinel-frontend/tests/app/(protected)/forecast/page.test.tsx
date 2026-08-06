import { render, screen } from "@testing-library/react";
import ForecastPage from "@/app/(protected)/forecast/page";
import { describe, it, expect } from "@jest/globals";

jest.mock("@/context/AuthContext", () => ({ useAuth: () => ({ user: { role: "admin" }, isAuthenticated: true }) }));

describe("ForecastPage", () => {
  it("should render", () => {
    render(<ForecastPage />);
    expect(document.body).toBeInTheDocument();
  });
});
