import { render, screen } from "@testing-library/react"
import { describe, it, expect } from "vitest"
import Avatar from "@/components/Avatar"

describe("Avatar", () => {
  it("renders successfully", () => {
    render(<Avatar name="Alice" />)
    expect(screen.getByRole("img", { name: /avatar/i })).toBeDefined()
  })

  it("displays the first letter for a regular name", () => {
    render(<Avatar name="alice" />)
    expect(screen.getByText("A")).toBeDefined()
  })

  it("displays the first two uppercase letters for a PascalCase name", () => {
    render(<Avatar name="JohnDoe" />)
    expect(screen.getByText("JD")).toBeDefined()
  })
})
