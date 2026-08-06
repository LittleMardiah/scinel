import { render, screen } from "@testing-library/react";
import { Footer } from "@/components/layout/Footer";
import { describe, it, expect } from "@jest/globals";

describe("Footer", () => {
  it("should render", () => {
    render(<Footer />);
    expect(document.body).toBeInTheDocument();
  });
});
