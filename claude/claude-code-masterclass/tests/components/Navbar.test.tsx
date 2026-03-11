import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi, beforeEach } from "vitest";

import Navbar from "@/components/Navbar";
import { useUser } from "@/context/AuthContext";

const mockSignOut = vi.fn();

vi.mock("firebase/auth", () => ({
  signOut: (...args: unknown[]) => mockSignOut(...args),
}));

vi.mock("@/lib/firebase/config", () => ({
  auth: {},
}));

vi.mock("@/context/AuthContext", () => ({
  useUser: vi.fn(),
}));

const mockUseUser = vi.mocked(useUser);

beforeEach(() => {
  mockUseUser.mockReturnValue({ user: null, isLoading: false });
  mockSignOut.mockReset();
});

describe("Navbar", () => {
  it("renders the main heading", () => {
    render(<Navbar />);

    const heading = screen.getByRole("heading", { level: 1 });
    expect(heading).toBeInTheDocument();
  });

  it("renders the Create Heist link", () => {
    render(<Navbar />);

    const createLink = screen.getByRole("link", { name: /create heist/i });
    expect(createLink).toBeInTheDocument();
    expect(createLink).toHaveAttribute("href", "/heists/create");
  });

  it("renders logout button when user is authenticated", () => {
    mockUseUser.mockReturnValue({
      user: { displayName: "Agent Fox" } as never,
      isLoading: false,
    });

    render(<Navbar />);

    expect(screen.getByRole("button", { name: /logout/i })).toBeInTheDocument();
  });

  it("does not render logout button when not authenticated", () => {
    render(<Navbar />);

    expect(
      screen.queryByRole("button", { name: /logout/i }),
    ).not.toBeInTheDocument();
  });

  it("calls signOut when logout button is clicked", async () => {
    mockUseUser.mockReturnValue({
      user: { displayName: "Agent Fox" } as never,
      isLoading: false,
    });
    mockSignOut.mockResolvedValue(undefined);

    render(<Navbar />);

    await userEvent.click(screen.getByRole("button", { name: /logout/i }));

    expect(mockSignOut).toHaveBeenCalledOnce();
  });
});
