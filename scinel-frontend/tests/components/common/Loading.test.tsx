import { render, screen } from "@testing-library/react";
import { Loading } from "@/components/common/Loading";
import { describe, it, expect } from "@jest/globals";

describe("Loading", () => {
  it("should render", () => {
    render(<Loading />);
    expect(document.body).toBeInTheDocument();
  });
});
