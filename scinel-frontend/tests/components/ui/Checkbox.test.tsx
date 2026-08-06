import { render, screen } from "@testing-library/react";
import { Checkbox } from "@/components/ui/Checkbox";
import { describe, it, expect } from "@jest/globals";

describe("Checkbox", () => {
  it("should render", () => {
    render(<Checkbox />);
    expect(document.body).toBeInTheDocument();
  });
});
