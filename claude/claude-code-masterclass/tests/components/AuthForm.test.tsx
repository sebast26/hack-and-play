import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi, beforeEach } from "vitest";

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

vi.mock("@/lib/firebase/signup", () => ({
  signUpUser: vi.fn().mockResolvedValue(undefined),
}));

const mockLoginUser = vi.fn();
vi.mock("@/lib/firebase/login", () => ({
  loginUser: (...args: unknown[]) => mockLoginUser(...args),
}));

import AuthForm from "@/components/AuthForm";

beforeEach(() => {
  mockLoginUser.mockReset();
});

describe("AuthForm", () => {
  describe("login mode", () => {
    it("renders email and password inputs", () => {
      render(<AuthForm mode="login" />);
      expect(screen.getByPlaceholderText("Email")).toBeDefined();
      expect(screen.getByPlaceholderText("Password")).toBeDefined();
    });

    it("renders a Login submit button", () => {
      render(<AuthForm mode="login" />);
      expect(screen.getByRole("button", { name: "Login" })).toBeDefined();
    });

    it("links to the signup page", () => {
      render(<AuthForm mode="login" />);
      const link = screen.getByRole("link", { name: /sign up/i });
      expect(link.getAttribute("href")).toBe("/signup");
    });

    it("shows a success message after valid credentials are submitted", async () => {
      mockLoginUser.mockResolvedValue(undefined);
      const user = userEvent.setup();
      render(<AuthForm mode="login" />);

      await user.type(screen.getByPlaceholderText("Email"), "test@example.com");
      await user.type(screen.getByPlaceholderText("Password"), "secret123");
      await user.click(screen.getByRole("button", { name: "Login" }));

      await waitFor(() => {
        expect(screen.getByText(/you're logged in/i)).toBeInTheDocument();
      });
    });

    it("does not show a success message when login fails", async () => {
      mockLoginUser.mockRejectedValue(new Error("Wrong password"));
      const user = userEvent.setup();
      render(<AuthForm mode="login" />);

      await user.type(screen.getByPlaceholderText("Email"), "test@example.com");
      await user.type(screen.getByPlaceholderText("Password"), "wrongpass");
      await user.click(screen.getByRole("button", { name: "Login" }));

      await waitFor(() => {
        expect(screen.queryByText(/you're logged in/i)).not.toBeInTheDocument();
      });
    });

    it("shows an error message when loginUser rejects", async () => {
      mockLoginUser.mockRejectedValue(new Error("Invalid credentials"));
      const user = userEvent.setup();
      render(<AuthForm mode="login" />);

      await user.type(screen.getByPlaceholderText("Email"), "test@example.com");
      await user.type(screen.getByPlaceholderText("Password"), "wrongpass");
      await user.click(screen.getByRole("button", { name: "Login" }));

      await waitFor(() => {
        expect(screen.getByText("Invalid credentials")).toBeInTheDocument();
      });
    });

    it("disables the submit button while loading", async () => {
      let resolve: () => void;
      mockLoginUser.mockReturnValue(
        new Promise<void>((res) => {
          resolve = res;
        }),
      );
      const user = userEvent.setup();
      render(<AuthForm mode="login" />);

      await user.type(screen.getByPlaceholderText("Email"), "test@example.com");
      await user.type(screen.getByPlaceholderText("Password"), "secret123");
      await user.click(screen.getByRole("button", { name: "Login" }));

      expect(screen.getByRole("button", { name: "Login" })).toBeDisabled();
      resolve!();
    });

    it("clears a previous error when the form is resubmitted", async () => {
      mockLoginUser.mockRejectedValueOnce(new Error("Wrong password"));
      mockLoginUser.mockResolvedValueOnce(undefined);
      const user = userEvent.setup();
      render(<AuthForm mode="login" />);

      await user.type(screen.getByPlaceholderText("Email"), "test@example.com");
      await user.type(screen.getByPlaceholderText("Password"), "wrongpass");
      await user.click(screen.getByRole("button", { name: "Login" }));

      await waitFor(() => {
        expect(screen.getByText("Wrong password")).toBeInTheDocument();
      });

      await user.click(screen.getByRole("button", { name: "Login" }));

      await waitFor(() => {
        expect(screen.queryByText("Wrong password")).not.toBeInTheDocument();
      });
    });
  });

  describe("signup mode", () => {
    it("renders email and password inputs", () => {
      render(<AuthForm mode="signup" />);
      expect(screen.getByPlaceholderText("Email")).toBeDefined();
      expect(screen.getByPlaceholderText("Password")).toBeDefined();
    });

    it("renders a Sign Up submit button", () => {
      render(<AuthForm mode="signup" />);
      expect(screen.getByRole("button", { name: "Sign Up" })).toBeDefined();
    });

    it("links to the login page", () => {
      render(<AuthForm mode="signup" />);
      const link = screen.getByRole("link", { name: /log in/i });
      expect(link.getAttribute("href")).toBe("/login");
    });
  });

  it("toggles password visibility", async () => {
    const user = userEvent.setup();
    render(<AuthForm mode="login" />);

    const passwordInput = screen.getByPlaceholderText("Password");
    expect(passwordInput.getAttribute("type")).toBe("password");

    await user.click(screen.getByRole("button", { name: /show password/i }));
    expect(passwordInput.getAttribute("type")).toBe("text");

    await user.click(screen.getByRole("button", { name: /hide password/i }));
    expect(passwordInput.getAttribute("type")).toBe("password");
  });
});
