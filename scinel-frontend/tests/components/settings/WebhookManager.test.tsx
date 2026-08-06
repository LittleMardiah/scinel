import { render, screen } from "@testing-library/react";
import { WebhookManager } from "@/components/settings/WebhookManager";
import { describe, it, expect } from "@jest/globals";

describe("WebhookManager", () => {
  it("should render", () => {
    render(<WebhookManager webhooks={[]} onUpdate={() => {}} />);
    expect(document.body).toBeInTheDocument();
  });
});
