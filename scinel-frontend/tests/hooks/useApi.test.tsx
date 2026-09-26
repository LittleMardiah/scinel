import { renderHook, waitFor, act } from "@testing-library/react";
import { useApi } from "@/hooks/useApi";
import { apiClient } from "@/lib/api";
import { describe, it, expect, jest, beforeEach, afterEach } from "@jest/globals";

describe("useApi", () => {
  let getSpy: ReturnType<typeof jest.spyOn>;

  beforeEach(() => {
    // Spy pada apiClient.get dan mock implementasi
    getSpy = jest.spyOn(apiClient, "get").mockImplementation(
      () => Promise.resolve({ data: { data: {} } })
    );
    jest.useFakeTimers();
  });

  afterEach(() => {
    getSpy.mockRestore();
    jest.useRealTimers();
  });

  it("should fetch data successfully", async () => {
    const mockData = { data: "test data" };
    getSpy.mockResolvedValueOnce({ data: { data: mockData } });

    const { result } = renderHook(() => useApi("/test", { autoFetch: true }));

    // Tunggu microtasks selesai
    await act(async () => {
      await Promise.resolve();
      jest.runAllTimers();
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.data).toEqual(mockData);
    });
  });

  it("should handle error", async () => {
    const mockError = new Error("Network error");
    getSpy.mockRejectedValueOnce(mockError);

    const { result } = renderHook(() => useApi("/test", { autoFetch: true }));

    await act(async () => {
      await Promise.resolve();
      jest.runAllTimers();
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toEqual(mockError);
    });
  });

  it("should call onSuccess callback", async () => {
    const onSuccess = jest.fn();
    const mockData = { data: "test" };
    getSpy.mockResolvedValueOnce({ data: { data: mockData } });

    renderHook(() => useApi("/test", { autoFetch: true, onSuccess }));

    await act(async () => {
      await Promise.resolve();
      jest.runAllTimers();
    });

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalledWith(mockData);
    });
  });

  it("should support refetch", async () => {
    const mockData = { data: "test" };
    getSpy.mockResolvedValue({ data: { data: mockData } });

    const { result } = renderHook(() => useApi("/test", { autoFetch: false }));

    await act(async () => {
      await result.current.refetch();
    });

    await waitFor(() => {
      expect(result.current.data).toEqual(mockData);
    });
  });
});
