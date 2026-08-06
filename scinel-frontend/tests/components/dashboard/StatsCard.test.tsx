import { render, screen } from "@testing-library/react";
import { StatsCard } from "@/components/dashboard/StatsCard";
import { describe, it, expect } from "@jest/globals";

describe("StatsCard", () => {
  it("should render title and value", () => {
    render(<StatsCard title="Total SKU" value="1,234" />);
    expect(screen.getByText("Total SKU")).toBeInTheDocument();
    expect(screen.getByText("1,234")).toBeInTheDocument();
  });

  it("should render icon when provided", () => {
    const { container } = render(
      <StatsCard title="Test" value="100" icon={<svg data-testid="icon" />} />
    );
    // Perbaiki selector: pakai kutip
    expect(container.querySelector('[data-testid="icon"]')).toBeInTheDocument();
  });
});