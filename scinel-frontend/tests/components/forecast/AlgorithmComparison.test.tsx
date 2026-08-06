import { render, screen } from "@testing-library/react";
import { AlgorithmComparison } from "@/components/forecast/AlgorithmComparison";
import { describe, it, expect } from "@jest/globals";

describe("AlgorithmComparison", () => {
  it("should render", () => {
    render(<AlgorithmComparison etsMape={10} arimaMape={12} />);
    expect(document.body).toBeInTheDocument();
  });
});
