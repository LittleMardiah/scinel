import * as Common from "@/components/common";
import { describe, it, expect } from "@jest/globals";

describe("Common index", () => {
  it("should export components", () => {
    expect(Common.Button).toBeDefined();
    expect(Common.Card).toBeDefined();
    expect(Common.Input).toBeDefined();
    expect(Common.Modal).toBeDefined();
  });
});
