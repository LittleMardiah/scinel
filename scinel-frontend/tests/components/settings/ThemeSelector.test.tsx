import { render, screen } from "@testing-library/react";
import { ThemeSelector } from "@/components/settings/ThemeSelector";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect } from "@jest/globals";

describe("ThemeSelector", () => {
  it("should render theme options", () => {
    render(
      <ThemeProvider>
        <ThemeSelector />
      </ThemeProvider>
    );
    // Cari label text yang muncul di DOM
    expect(screen.getByText(/☀️ Terang/i)).toBeInTheDocument();
    expect(screen.getByText(/🌙 Gelap/i)).toBeInTheDocument();
    expect(screen.getByText(/🖥️ Sistem/i)).toBeInTheDocument();
  });
});
