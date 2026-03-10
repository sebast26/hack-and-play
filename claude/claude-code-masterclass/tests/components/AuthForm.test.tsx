import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import AuthForm from "@/components/AuthForm";

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

  it("logs email and password on submit", async () => {
    const user = userEvent.setup();
    const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});

    render(<AuthForm mode="login" />);

    await user.type(screen.getByPlaceholderText("Email"), "test@example.com");
    await user.type(screen.getByPlaceholderText("Password"), "secret123");
    await user.click(screen.getByRole("button", { name: "Login" }));

    expect(consoleSpy).toHaveBeenCalledWith({
      email: "test@example.com",
      password: "secret123",
    });

    consoleSpy.mockRestore();
  });
});
