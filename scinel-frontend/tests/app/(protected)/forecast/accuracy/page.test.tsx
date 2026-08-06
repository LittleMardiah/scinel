import { render, screen } from "@testing-library/react";
import ForecastAccuracyPage from "@/app/(protected)/forecast/accuracy/page";
import { describe, it, expect } from "@jest/globals";

jest.mock("@/context/AuthContext", () => ({ useAuth: () => ({ user: { role: "admin" }, isAuthenticated: true }) }));

describe("ForecastAccuracyPage", () => {
  it("should render", () => {
    render(<ForecastAccuracyPage />);
    expect(document.body).toBeInTheDocument();
  });
});
