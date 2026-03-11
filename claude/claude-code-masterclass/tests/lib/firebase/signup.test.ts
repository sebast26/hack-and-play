import { describe, it, expect, vi, beforeEach } from "vitest";

const mockCreateUserWithEmailAndPassword = vi.fn();
const mockUpdateProfile = vi.fn();
const mockSetDoc = vi.fn();
const mockDoc = vi.fn();
const mockGenerateCodename = vi.fn(() => "SwiftFoxStrikes");

vi.mock("firebase/auth", () => ({
  createUserWithEmailAndPassword: (...args: unknown[]) =>
    mockCreateUserWithEmailAndPassword(...args),
  updateProfile: (...args: unknown[]) => mockUpdateProfile(...args),
}));

vi.mock("firebase/firestore", () => ({
  setDoc: (...args: unknown[]) => mockSetDoc(...args),
  doc: (...args: unknown[]) => mockDoc(...args),
}));

vi.mock("@/lib/firebase/config", () => ({
  auth: {},
  db: {},
}));

vi.mock("@/lib/codename", () => ({
  generateCodename: () => mockGenerateCodename(),
}));

import { signUpUser } from "@/lib/firebase/signup";

describe("signUpUser", () => {
  const fakeUser = { uid: "test-uid-123" };

  beforeEach(() => {
    vi.clearAllMocks();
    mockCreateUserWithEmailAndPassword.mockResolvedValue({ user: fakeUser });
    mockUpdateProfile.mockResolvedValue(undefined);
    mockSetDoc.mockResolvedValue(undefined);
    mockDoc.mockReturnValue("doc-ref");
  });

  it("calls createUserWithEmailAndPassword with provided credentials", async () => {
    await signUpUser("user@example.com", "password123");
    expect(mockCreateUserWithEmailAndPassword).toHaveBeenCalledWith(
      {},
      "user@example.com",
      "password123",
    );
  });

  it("calls updateProfile with the generated codename as displayName", async () => {
    await signUpUser("user@example.com", "password123");
    expect(mockUpdateProfile).toHaveBeenCalledWith(fakeUser, {
      displayName: "SwiftFoxStrikes",
    });
  });

  it("creates a Firestore document with id and codename but no email", async () => {
    await signUpUser("user@example.com", "password123");
    expect(mockSetDoc).toHaveBeenCalledWith("doc-ref", {
      id: "test-uid-123",
      codename: "SwiftFoxStrikes",
    });
    const docData = mockSetDoc.mock.calls[0][1];
    expect(docData).not.toHaveProperty("email");
  });

  it("uses the user UID as the Firestore document ID", async () => {
    await signUpUser("user@example.com", "password123");
    expect(mockDoc).toHaveBeenCalledWith({}, "users", "test-uid-123");
  });

  it("propagates errors from createUserWithEmailAndPassword", async () => {
    const error = new Error("Firebase: email already in use");
    mockCreateUserWithEmailAndPassword.mockRejectedValue(error);
    await expect(
      signUpUser("taken@example.com", "password123"),
    ).rejects.toThrow("Firebase: email already in use");
  });
});
