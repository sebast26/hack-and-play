import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import type { User } from "firebase/auth";

import { useUser } from "@/context/AuthContext";

const mockReplace = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace: mockReplace }),
}));

vi.mock("@/context/AuthContext", () => ({
  useUser: vi.fn(),
}));

vi.mock("lucide-react", () => ({
  Clock8: () => <svg data-testid="clock-icon" />,
}));

vi.mock("@/components/Navbar", () => ({
  default: () => <nav data-testid="navbar" />,
}));

const mockUseUser = vi.mocked(useUser);

import DashboardLayout from "@/app/(dashboard)/layout";

beforeEach(() => {
  vi.clearAllMocks();
});

describe("DashboardLayout", () => {
  it("renders loader when auth state is loading", () => {
    mockUseUser.mockReturnValue({ user: null, isLoading: true });

    render(<DashboardLayout>content</DashboardLayout>);

    expect(screen.getByTestId("clock-icon")).toBeInTheDocument();
    expect(screen.queryByText("content")).not.toBeInTheDocument();
  });

  it("renders loader and redirects to /login when user is unauthenticated", () => {
    mockUseUser.mockReturnValue({ user: null, isLoading: false });

    render(<DashboardLayout>content</DashboardLayout>);

    expect(screen.getByTestId("clock-icon")).toBeInTheDocument();
    expect(screen.queryByText("content")).not.toBeInTheDocument();
    expect(mockReplace).toHaveBeenCalledWith("/login");
  });

  it("renders Navbar and children when user is authenticated", () => {
    mockUseUser.mockReturnValue({
      user: { uid: "123" } as User,
      isLoading: false,
    });

    render(<DashboardLayout>content</DashboardLayout>);

    expect(screen.getByTestId("navbar")).toBeInTheDocument();
    expect(screen.getByText("content")).toBeInTheDocument();
    expect(screen.queryByTestId("clock-icon")).not.toBeInTheDocument();
    expect(mockReplace).not.toHaveBeenCalled();
  });
});
