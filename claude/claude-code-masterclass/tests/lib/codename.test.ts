import { describe, it, expect } from "vitest";
import { generateCodename } from "@/lib/codename";

describe("generateCodename", () => {
  it("returns a non-empty string", () => {
    expect(generateCodename()).toBeTruthy();
  });

  it("returns a valid PascalCase string composed of three words", () => {
    const codename = generateCodename();
    // Three PascalCase words: each starts with uppercase, rest lowercase, no separators
    expect(codename).toMatch(/^([A-Z][a-z]+){3}$/);
  });

  it("returns different values across multiple calls", () => {
    const results = new Set(
      Array.from({ length: 10 }, () => generateCodename()),
    );
    expect(results.size).toBeGreaterThan(1);
  });
});
