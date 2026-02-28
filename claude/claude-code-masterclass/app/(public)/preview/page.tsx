// preview page for newly created UI components

import Skeleton from "@/components/Skeleton"
import Avatar from "@/components/Avatar"

export default function PreviewPage() {
  return (
    <div className="page-content">
      <h2>Preview</h2>
      <div className="preview-grid">
        <Skeleton />
        <Skeleton />
        <Skeleton />
        <Skeleton />
        <Skeleton />
        <Skeleton />
      </div>

      <h2>Avatar</h2>
      <div className="preview-grid">
        <Avatar name="alice" />
        <Avatar name="JohnDoe" />
        <Avatar name="MaryJane" />
        <Avatar name="Bob" />
      </div>
    </div>
  )
}
