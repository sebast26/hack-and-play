import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";

// component imports
import Navbar from "@/components/Navbar";

vi.mock("@/context/AuthContext", () => ({
  useUser: () => ({ user: null, isLoading: false }),
}));

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
});
