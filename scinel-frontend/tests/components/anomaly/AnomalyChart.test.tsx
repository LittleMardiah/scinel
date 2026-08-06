import { render, screen } from "@testing-library/react";
import { AnomalyChart } from "@/components/anomaly/AnomalyChart";
import { describe, it, expect } from "@jest/globals";

describe("AnomalyChart", () => {
  it("should render", () => {
    render(<AnomalyChart data={[]} />);
    expect(document.body).toBeInTheDocument();
  });
});
