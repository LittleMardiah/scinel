import { render, screen } from "@testing-library/react";
import { ThemeToggle } from "@/components/ui/ThemeToggle";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

describe("ThemeToggle", () => {
  it("should render", () => {
    render(
      <ThemeProvider>
        <ThemeToggle />
      </ThemeProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
