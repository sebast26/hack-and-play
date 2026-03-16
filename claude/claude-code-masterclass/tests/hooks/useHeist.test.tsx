import { renderHook, act, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import type { Heist } from "@/types/firestore";

vi.mock("@/lib/firebase/config", () => ({ db: {} }));

const mockOnSnapshot = vi.fn();
vi.mock("firebase/firestore", () => ({
  collection: vi.fn(() => ({ withConverter: vi.fn(() => ({})) })),
  query: vi.fn(),
  where: vi.fn(),
  onSnapshot: (...args: unknown[]) => mockOnSnapshot(...args),
}));

const mockUser = { uid: "u1", displayName: "Viper" };
let mockUserReturn: { user: typeof mockUser | null } = { user: mockUser };

vi.mock("@/context/AuthContext", () => ({
  useUser: () => mockUserReturn,
}));

import { useHeist } from "@/hooks/useHeist";

function makeSnapshot(docs: Partial<Heist>[]) {
  return {
    docs: docs.map((d) => ({ data: () => d })),
  };
}

function setupSnapshot(docs: Partial<Heist>[]) {
  const unsubscribe = vi.fn();
  mockOnSnapshot.mockImplementation(
    (_q: unknown, onNext: (snap: unknown) => void) => {
      onNext(makeSnapshot(docs));
      return unsubscribe;
    },
  );
  return unsubscribe;
}

const future = new Date(Date.now() + 99999999);
const past = new Date(Date.now() - 99999999);

beforeEach(() => {
  vi.clearAllMocks();
  mockUserReturn = { user: mockUser };
});

describe("useHeist", () => {
  it("returns loading:true and empty array before snapshot fires", () => {
    mockOnSnapshot.mockImplementation(() => vi.fn()); // never calls onNext
    const { result } = renderHook(() => useHeist("active"));
    expect(result.current.loading).toBe(true);
    expect(result.current.heists).toEqual([]);
  });

  it("returns active heists assigned to the current user with a future deadline", async () => {
    const activeHeist: Partial<Heist> = {
      id: "h1",
      title: "Steal the stapler",
      assignedTo: "u1",
      createdBy: "u2",
      deadline: future,
      finalStatus: null,
    };
    setupSnapshot([activeHeist]);

    const { result } = renderHook(() => useHeist("active"));

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.heists).toHaveLength(1);
    expect(result.current.heists[0].title).toBe("Steal the stapler");
  });

  it("returns assigned heists created by the current user with a future deadline", async () => {
    const assignedHeist: Partial<Heist> = {
      id: "h2",
      title: "Plant the rubber duck",
      createdBy: "u1",
      assignedTo: "u2",
      deadline: future,
      finalStatus: null,
    };
    setupSnapshot([assignedHeist]);

    const { result } = renderHook(() => useHeist("assigned"));

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.heists).toHaveLength(1);
    expect(result.current.heists[0].title).toBe("Plant the rubber duck");
  });

  it("filters expired heists client-side to only those with a non-null finalStatus", async () => {
    const withStatus: Partial<Heist> = {
      id: "h3",
      title: "Done",
      deadline: past,
      finalStatus: "success",
    };
    const withoutStatus: Partial<Heist> = {
      id: "h4",
      title: "Pending",
      deadline: past,
      finalStatus: null,
    };
    setupSnapshot([withStatus, withoutStatus]);

    const { result } = renderHook(() => useHeist("expired"));

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.heists).toHaveLength(1);
    expect(result.current.heists[0].title).toBe("Done");
  });

  it("does not call onSnapshot when user is null", () => {
    mockUserReturn = { user: null };
    renderHook(() => useHeist("active"));
    expect(mockOnSnapshot).not.toHaveBeenCalled();
  });

  it("calls the unsubscribe function on unmount", async () => {
    const unsubscribe = setupSnapshot([]);
    const { unmount } = renderHook(() => useHeist("active"));
    await waitFor(() => expect(mockOnSnapshot).toHaveBeenCalled());
    unmount();
    expect(unsubscribe).toHaveBeenCalled();
  });

  it("sets error when onSnapshot fires its error callback", async () => {
    mockOnSnapshot.mockImplementation(
      (_q: unknown, _onNext: unknown, onError: (e: Error) => void) => {
        onError(new Error("Permission denied"));
        return vi.fn();
      },
    );

    const { result } = renderHook(() => useHeist("active"));

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.error).toBe("Permission denied");
    expect(result.current.heists).toEqual([]);
  });
});
