import { render, screen } from "@testing-library/react";
import AnomalyAnalyticsPage from "@/app/(protected)/anomalies/analytics/page";
import { describe, it, expect } from "@jest/globals";

jest.mock("@/context/AuthContext", () => ({ useAuth: () => ({ user: { role: "admin" }, isAuthenticated: true }) }));

describe("AnomalyAnalyticsPage", () => {
  it("should render", () => {
    render(<AnomalyAnalyticsPage />);
    expect(document.body).toBeInTheDocument();
  });
});
