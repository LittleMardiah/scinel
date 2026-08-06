import { render, screen } from "@testing-library/react";
import { MobileNav } from "@/components/layout/MobileNav";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import { describe, it, expect, jest } from "@jest/globals";

jest.mock("next/navigation", () => ({
  usePathname: () => "/",
  useRouter: () => ({ push: jest.fn() }),
}));

describe("MobileNav", () => {
  it("should render", () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <MobileNav items={[]} />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(document.body).toBeInTheDocument();
  });
});
