import { render, screen } from "@testing-library/react";
import { AnomalyCard } from "@/components/anomaly/AnomalyCard";
import { describe, it, expect } from "@jest/globals";

describe("AnomalyCard", () => {
  it("should render", () => {
    render(<AnomalyCard title="Test" severity="high" />);
    expect(document.body).toBeInTheDocument();
  });
});
