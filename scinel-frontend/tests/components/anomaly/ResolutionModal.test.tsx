import { render, screen } from "@testing-library/react";
import { ResolutionModal } from "@/components/anomaly/ResolutionModal";
import { describe, it, expect } from "@jest/globals";

describe("ResolutionModal", () => {
  it("should render", () => {
    render(<ResolutionModal isOpen={true} onClose={() => {}} onResolve={() => {}} />);
    expect(document.body).toBeInTheDocument();
  });
});
