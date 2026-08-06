import { render, screen } from "@testing-library/react";
import { TextArea } from "@/components/ui/TextArea";
import { describe, it, expect } from "@jest/globals";

describe("TextArea", () => {
  it("should render", () => {
    render(<TextArea />);
    expect(document.body).toBeInTheDocument();
  });
});
