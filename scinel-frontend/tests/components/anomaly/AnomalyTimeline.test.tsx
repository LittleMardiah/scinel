import { render, screen } from "@testing-library/react";
import { AnomalyTimeline } from "@/components/anomaly/AnomalyTimeline";
import { describe, it, expect } from "@jest/globals";

describe("AnomalyTimeline", () => {
  it("should render", () => {
    render(<AnomalyTimeline events={[]} />);
    expect(document.body).toBeInTheDocument();
  });
});
