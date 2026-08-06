import { renderHook, waitFor } from "@testing-library/react";
import { useDashboardData } from "@/hooks/useDashboardData";
import { describe, it, expect, jest } from "@jest/globals";

jest.mock("@/lib/api", () => ({
  apiClient: {
    get: jest.fn().mockResolvedValue({
      data: {
        data: {
          total_sku: 100,
          health: { status: "online" },
        },
      },
    }),
  },
}));

describe("useDashboardData", () => {
  it("should fetch data", async () => {
    const { result } = renderHook(() => useDashboardData());
    // isLoading awalnya true, tapi cepat berubah. Gunakan waitFor untuk assertion yang stabil.
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });
    expect(result.current.data).toBeDefined();
  });
});
