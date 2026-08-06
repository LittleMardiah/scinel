import { render, screen } from "@testing-library/react";
import { DatePicker } from "@/components/ui/DatePicker";
import { describe, it, expect } from "@jest/globals";

describe("DatePicker", () => {
  it("should render", () => {
    render(<DatePicker />);
    expect(document.body).toBeInTheDocument();
  });
});
