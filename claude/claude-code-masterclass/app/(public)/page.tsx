// this page should be used only as a splash page to decide where a user should be navigated to
// when logged in --> to /heists
// when not logged in --> to /login

import { Clock8 } from "lucide-react"

export default function Home() {
  return (
    <div className="center-content">
      <div className="page-content">
        <h1>
          P<Clock8 className="logo" strokeWidth={2.75} />cket Heist
        </h1>
        <div>Tiny missions. Big office mischief.</div>
        <p>
          Welcome to Pocket Heist — the app that turns your office into a playground.
          Plan covert missions, assign tasks to unsuspecting colleagues, and track the
          chaos as it unfolds. No heist too small, no mischief too grand.
        </p>
      </div>
    </div>
  )
}
