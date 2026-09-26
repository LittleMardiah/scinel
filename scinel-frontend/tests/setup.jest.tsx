import "@testing-library/jest-dom";
import "@testing-library/jest-dom/jest-globals";
import { jest } from "@jest/globals";

// ===== MOCK GLOBAL NEXT/NAVIGATION =====
jest.mock("next/navigation", () => ({
  redirect: jest.fn(),
  permanentRedirect: jest.fn(),
  notFound: jest.fn(),
  useRouter: jest.fn(() => ({
    push: jest.fn(),
    replace: jest.fn(),
    prefetch: jest.fn(),
    back: jest.fn(),
    forward: jest.fn(),
    refresh: jest.fn(),
    pathname: "/",
    query: {},
  })),
  usePathname: jest.fn(() => "/"),
  useSearchParams: jest.fn(() => new URLSearchParams()),
  useParams: jest.fn(() => ({})),
}));

// ===== MOCK GLOBAL NEXT-THEMES =====
jest.mock("next-themes", () => ({
  useTheme: jest.fn(() => ({
    theme: "light",
    setTheme: jest.fn(),
    resolvedTheme: "light",
    themes: ["light", "dark", "system"],
  })),
  ThemeProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

// ===== MOCK WINDOW.MATCHMEDIA =====
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: jest.fn().mockImplementation((query: unknown) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(),
    removeListener: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// ===== MOCK INTERSECTIONOBSERVER =====
global.IntersectionObserver = jest.fn(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
})) as unknown as typeof IntersectionObserver;

// ===== MOCK RECHARTS =====
jest.mock("recharts", () => {
  const React = require("react");
  return {
    ResponsiveContainer: ({ children }: { children: React.ReactNode }) => (
      <div style={{ width: "100%", height: "100%" }}>{children}</div>
    ),
    ComposedChart: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    LineChart: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    BarChart: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    PieChart: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    Line: () => <div />,
    Area: () => <div />,
    Bar: () => <div />,
    Pie: () => <div />,
    Cell: () => <div />,
    XAxis: () => <div />,
    YAxis: () => <div />,
    CartesianGrid: () => <div />,
    Tooltip: () => <div />,
    Legend: () => <div />,
  };
});

// ===== MOCK THEMETOGGLE =====
jest.mock("@/components/ui/ThemeToggle", () => ({
  ThemeToggle: () => <button data-testid="theme-toggle">Toggle Theme</button>,
}));

// ===== TIDAK ADA MOCK UNTUK LAYOUT COMPONENTS DI SINI =====
// (Mock dilakukan di masing-masing test file jika diperlukan)

// console.error = jest.fn();