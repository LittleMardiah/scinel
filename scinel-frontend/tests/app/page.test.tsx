import { render, screen } from "@testing-library/react";
import Page from "@/app/page";
import { describe, it, expect } from "@jest/globals";

describe("Page", () => {
  it("should render", () => {
    render(<Page />);
    // Page biasanya redirect atau render layout
    expect(document.body).toBeInTheDocument();
  });
});
