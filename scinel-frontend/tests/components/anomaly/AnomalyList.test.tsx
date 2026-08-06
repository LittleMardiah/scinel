import { render, screen } from "@testing-library/react";
import { AnomalyList } from "@/components/anomaly/AnomalyList";
import { describe, it, expect, jest } from "@jest/globals";

const mockData = [
  {
    id: "1",
    sku_code: "SKU-001",
    type: "unusual_drop",
    severity: "high" as const,
    description: "Unusual drop",
    z_score: 3.5,
    status: "open" as const,
    detected_at: "2026-07-29T07:00:00Z",
  },
  {
    id: "2",
    sku_code: "SKU-002",
    type: "minor_spike",
    severity: "low" as const,
    description: "Minor spike",
    z_score: 1.2,
    status: "resolved" as const,
    detected_at: "2026-07-28T07:00:00Z",
  },
];

describe("AnomalyList", () => {
  it("should render anomaly items", () => {
    render(
      <AnomalyList
        data={mockData}
        onSort={jest.fn()}
        sortKey="status"
        sortOrder="asc"
        onAcknowledge={jest.fn()}
        onResolve={jest.fn()}
      />
    );
    expect(screen.getByText(/SKU-001/i)).toBeInTheDocument();
    expect(screen.getByText(/SKU-002/i)).toBeInTheDocument();
    expect(screen.getByText(/Terbuka/i)).toBeInTheDocument();
    expect(screen.getByText(/Selesai/i)).toBeInTheDocument();
  });
});
