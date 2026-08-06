import { render, screen } from "@testing-library/react";
import { AccuracyMetrics } from "@/components/forecast/AccuracyMetrics";
import { describe, it, expect } from "@jest/globals";

describe("AccuracyMetrics", () => {
  it("should render", () => {
    render(<AccuracyMetrics mape={10} rmse={5} mae={3} />);
    expect(document.body).toBeInTheDocument();
  });
});
