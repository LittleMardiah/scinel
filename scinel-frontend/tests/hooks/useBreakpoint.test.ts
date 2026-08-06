import { renderHook } from "@testing-library/react";
import { useBreakpoint } from "@/hooks/useBreakpoint";
import { describe, it, expect } from "@jest/globals";

describe("useBreakpoint", () => {
  it("should return breakpoint", () => {
    const { result } = renderHook(() => useBreakpoint());
    expect(result.current).toBeDefined();
  });
});
