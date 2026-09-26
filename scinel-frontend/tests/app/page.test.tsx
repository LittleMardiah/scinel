import { redirect } from "next/navigation";
import Page from "@/app/page";
import { describe, it, expect } from "@jest/globals";

describe("Page", () => {
  it("should redirect to /login", () => {
    Page();
    expect(redirect).toHaveBeenCalledWith("/login");
  });
});
