import { apiClient, getApiKeys, createApiKey, revokeApiKey, rotateApiKey, getWebhooks, createWebhook, deleteWebhook, toggleWebhook } from "@/lib/api";
import { describe, it, expect } from "@jest/globals";

describe("API Library", () => {
  it("should export apiClient", () => {
    expect(apiClient).toBeDefined();
  });

  it("should export API key functions", () => {
    expect(getApiKeys).toBeDefined();
    expect(createApiKey).toBeDefined();
    expect(revokeApiKey).toBeDefined();
    expect(rotateApiKey).toBeDefined();
  });

  it("should export webhook functions", () => {
    expect(getWebhooks).toBeDefined();
    expect(createWebhook).toBeDefined();
    expect(deleteWebhook).toBeDefined();
    expect(toggleWebhook).toBeDefined();
  });
});
