import { render, screen } from "@testing-library/react";
import { KPIGrid } from "@/components/dashboard/KPIGrid";
import { describe, it, expect } from "@jest/globals";

describe("KPIGrid", () => {
  it("should render component", () => {
    render(<KPIGrid />);
    // Cari salah satu elemen yang pasti ada (pakai getAllByText untuk handle duplikat)
    const elements = screen.getAllByText(/KPIGrid/i);
    expect(elements.length).toBeGreaterThan(0);
  });
});