import { render, screen } from "@testing-library/react";
import { Breadcrumb } from "@/components/common/Breadcrumb";
import { describe, it, expect } from "@jest/globals";

describe("Breadcrumb", () => {
  it("should render", () => {
    render(<Breadcrumb items={[]} />);
    expect(document.body).toBeInTheDocument();
  });
});
