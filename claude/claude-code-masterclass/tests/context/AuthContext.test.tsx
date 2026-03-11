import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import type { User } from "firebase/auth";

// Mock firebase/auth before importing AuthContext
const mockOnAuthStateChanged = vi.fn();

vi.mock("firebase/auth", () => ({
  onAuthStateChanged: (...args: unknown[]) => mockOnAuthStateChanged(...args),
  getAuth: vi.fn(),
}));

vi.mock("@/lib/firebase/config", () => ({
  auth: {},
  default: {},
}));

import { AuthProvider, useUser } from "@/context/AuthContext";
import { renderHook } from "@testing-library/react";

function renderWithProvider(ui: React.ReactNode) {
  return render(<AuthProvider>{ui}</AuthProvider>);
}

function setupAuthListener(firebaseUser: User | null) {
  mockOnAuthStateChanged.mockImplementation(
    (_auth: unknown, callback: (user: User | null) => void) => {
      callback(firebaseUser);
      return vi.fn(); // unsubscribe
    },
  );
}

beforeEach(() => {
  vi.clearAllMocks();
});

describe("AuthProvider", () => {
  it("renders children without errors", () => {
    setupAuthListener(null);

    renderWithProvider(<div>child content</div>);

    expect(screen.getByText("child content")).toBeInTheDocument();
  });
});

describe("useUser", () => {
  it("returns null user and isLoading false when not authenticated", async () => {
    setupAuthListener(null);

    const { result } = renderHook(() => useUser(), {
      wrapper: AuthProvider,
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.user).toBeNull();
  });

  it("returns the user object when authenticated", async () => {
    const mockUser = {
      uid: "123",
      email: "test@example.com",
      displayName: "Test User",
    } as User;
    setupAuthListener(mockUser);

    const { result } = renderHook(() => useUser(), {
      wrapper: AuthProvider,
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.user).toEqual(mockUser);
  });

  it("throws when called outside AuthProvider", () => {
    expect(() => {
      renderHook(() => useUser());
    }).toThrow("useUser must be used within AuthProvider");
  });
});
