import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi, beforeEach } from "vitest";

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

vi.mock("@/context/AuthContext", () => ({
  useUser: () => ({ user: { uid: "u1", displayName: "Viper" } }),
}));

vi.mock("@/lib/firebase/config", () => ({ db: {} }));

const mockAddDoc = vi.fn();
const mockGetDocs = vi.fn();
vi.mock("firebase/firestore", () => ({
  collection: vi.fn(),
  getDocs: (...args: unknown[]) => mockGetDocs(...args),
  addDoc: (...args: unknown[]) => mockAddDoc(...args),
  serverTimestamp: vi.fn(() => "SERVER_TIMESTAMP"),
}));

const mockUsers = [
  { id: "u1", codename: "Viper" },
  { id: "u2", codename: "Ghost" },
];

function makeMockSnapshot(docs: { id: string; codename: string }[]) {
  return {
    docs: docs.map((d) => ({ data: () => d })),
  };
}

import CreateHeistForm from "@/components/CreateHeistForm";

beforeEach(() => {
  mockPush.mockReset();
  mockAddDoc.mockReset();
  mockGetDocs.mockResolvedValue(makeMockSnapshot(mockUsers));
});

describe("CreateHeistForm", () => {
  it("renders title input, description textarea, and assignee dropdown", async () => {
    render(<CreateHeistForm />);
    expect(screen.getByPlaceholderText("Mission title")).toBeDefined();
    expect(
      screen.getByPlaceholderText("Describe the mission..."),
    ).toBeDefined();
    await waitFor(() => {
      expect(screen.getByRole("combobox")).toBeDefined();
    });
  });

  it("populates the dropdown with users from Firestore", async () => {
    render(<CreateHeistForm />);
    await waitFor(() => {
      expect(screen.getByRole("option", { name: "Viper" })).toBeDefined();
      expect(screen.getByRole("option", { name: "Ghost" })).toBeDefined();
    });
  });

  it("calls addDoc with a CreateHeistInput-shaped payload on submit", async () => {
    mockAddDoc.mockResolvedValue({ id: "new-heist" });
    const user = userEvent.setup();
    render(<CreateHeistForm />);

    await waitFor(() => screen.getByRole("option", { name: "Ghost" }));

    await user.type(
      screen.getByPlaceholderText("Mission title"),
      "Office Takeover",
    );
    await user.type(
      screen.getByPlaceholderText("Describe the mission..."),
      "Steal the stapler.",
    );
    await user.selectOptions(screen.getByRole("combobox"), "Ghost");
    await user.click(screen.getByRole("button", { name: /launch heist/i }));

    await waitFor(() => {
      expect(mockAddDoc).toHaveBeenCalledOnce();
      const payload = mockAddDoc.mock.calls[0][1];
      expect(payload.title).toBe("Office Takeover");
      expect(payload.description).toBe("Steal the stapler.");
      expect(payload.createdBy).toBe("u1");
      expect(payload.createdByCodename).toBe("Viper");
      expect(payload.assignedTo).toBe("u2");
      expect(payload.assignedToCodename).toBe("Ghost");
      expect(payload.createdAt).toBe("SERVER_TIMESTAMP");
      expect(payload.finalStatus).toBeNull();
      expect(payload.deadline).toBeInstanceOf(Date);
    });
  });

  it("redirects to /heists after successful submit", async () => {
    mockAddDoc.mockResolvedValue({ id: "new-heist" });
    const user = userEvent.setup();
    render(<CreateHeistForm />);

    await waitFor(() => screen.getByRole("option", { name: "Ghost" }));

    await user.type(
      screen.getByPlaceholderText("Mission title"),
      "Office Takeover",
    );
    await user.type(
      screen.getByPlaceholderText("Describe the mission..."),
      "Steal the stapler.",
    );
    await user.selectOptions(screen.getByRole("combobox"), "Ghost");
    await user.click(screen.getByRole("button", { name: /launch heist/i }));

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith("/heists");
    });
  });

  it("disables the submit button while submission is in progress", async () => {
    let resolve: () => void;
    mockAddDoc.mockReturnValue(
      new Promise<void>((res) => {
        resolve = res;
      }),
    );
    const user = userEvent.setup();
    render(<CreateHeistForm />);

    await waitFor(() => screen.getByRole("option", { name: "Ghost" }));

    await user.type(
      screen.getByPlaceholderText("Mission title"),
      "Office Takeover",
    );
    await user.type(
      screen.getByPlaceholderText("Describe the mission..."),
      "Steal the stapler.",
    );
    await user.selectOptions(screen.getByRole("combobox"), "Ghost");
    await user.click(screen.getByRole("button", { name: /launch heist/i }));

    expect(screen.getByRole("button", { name: /launching/i })).toBeDisabled();
    resolve!();
  });

  it("shows an error message when addDoc rejects and keeps form interactive", async () => {
    mockAddDoc.mockRejectedValue(new Error("Permission denied"));
    const user = userEvent.setup();
    render(<CreateHeistForm />);

    await waitFor(() => screen.getByRole("option", { name: "Ghost" }));

    await user.type(
      screen.getByPlaceholderText("Mission title"),
      "Office Takeover",
    );
    await user.type(
      screen.getByPlaceholderText("Describe the mission..."),
      "Steal the stapler.",
    );
    await user.selectOptions(screen.getByRole("combobox"), "Ghost");
    await user.click(screen.getByRole("button", { name: /launch heist/i }));

    await waitFor(() => {
      expect(screen.getByText("Permission denied")).toBeDefined();
    });
    expect(
      screen.getByRole("button", { name: /launch heist/i }),
    ).not.toBeDisabled();
  });
});
