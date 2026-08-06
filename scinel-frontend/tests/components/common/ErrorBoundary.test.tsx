import { render, screen } from "@testing-library/react";
import { ErrorBoundary } from "@/components/common/ErrorBoundary";
import { describe, it, expect } from "@jest/globals";

describe("ErrorBoundary", () => {
  it("should render children", () => {
    render(
      <ErrorBoundary>
        <div>Test Child</div>
      </ErrorBoundary>
    );
    expect(screen.getByText("Test Child")).toBeInTheDocument();
  });
});
