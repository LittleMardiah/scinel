import { render, screen } from "@testing-library/react";
import { ForecastTable } from "@/components/forecast/ForecastTable";
import { describe, it, expect } from "@jest/globals";

describe("ForecastTable", () => {
  it("should render", () => {
    render(<ForecastTable data={[]} />);
    expect(document.body).toBeInTheDocument();
  });
});
