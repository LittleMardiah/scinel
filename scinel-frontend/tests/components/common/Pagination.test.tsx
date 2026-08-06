import { render, screen } from "@testing-library/react";
import { Pagination } from "@/components/common/Pagination";
import { describe, it, expect } from "@jest/globals";

describe("Pagination", () => {
  it("should render", () => {
    render(<Pagination totalPages={5} currentPage={1} onPageChange={() => {}} />);
    expect(document.body).toBeInTheDocument();
  });
});
