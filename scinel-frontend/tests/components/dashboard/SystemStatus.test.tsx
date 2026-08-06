import { render, screen } from "@testing-library/react";
import { SystemStatus } from "@/components/dashboard/SystemStatus";
import { describe, it, expect } from "@jest/globals";

describe("SystemStatus", () => {
  it("should show online status", () => {
    render(<SystemStatus status="online" />);
    expect(screen.getByText(/Online/i)).toBeInTheDocument();
  });

  it("should show offline status", () => {
    render(<SystemStatus status="offline" />);
    expect(screen.getByText(/Offline/i)).toBeInTheDocument();
  });
});
