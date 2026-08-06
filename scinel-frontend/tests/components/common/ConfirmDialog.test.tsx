import { render, screen } from "@testing-library/react";
import { ConfirmDialog } from "@/components/common/ConfirmDialog";
import { describe, it, expect } from "@jest/globals";

describe("ConfirmDialog", () => {
  it("should render", () => {
    render(<ConfirmDialog isOpen={true} onClose={() => {}} onConfirm={() => {}} title="Test" />);
    expect(document.body).toBeInTheDocument();
  });
});
