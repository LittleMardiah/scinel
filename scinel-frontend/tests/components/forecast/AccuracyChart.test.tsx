import { render, screen } from "@testing-library/react";
import { AccuracyChart } from "@/components/forecast/AccuracyChart";
import { describe, it, expect } from "@jest/globals";

describe("AccuracyChart", () => {
  it("should render", () => {
    render(<AccuracyChart data={[]} />);
    expect(document.body).toBeInTheDocument();
  });
});
