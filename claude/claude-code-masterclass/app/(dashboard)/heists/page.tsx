"use client";

import { useHeist } from "@/hooks/useHeist";

function HeistList({
  mode,
  heading,
}: {
  mode: "active" | "assigned" | "expired";
  heading: string;
}) {
  const { heists, loading } = useHeist(mode);

  return (
    <div>
      <h2>{heading}</h2>
      {loading ? (
        <p>Loading...</p>
      ) : heists.length === 0 ? (
        <p>No heists yet.</p>
      ) : (
        <ul>
          {heists.map((h) => (
            <li key={h.id}>{h.title}</li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default function HeistsPage() {
  return (
    <div className="page-content">
      <HeistList mode="active" heading="Your Active Heists" />
      <HeistList mode="assigned" heading="Heists You've Assigned" />
      <HeistList mode="expired" heading="All Expired Heists" />
    </div>
  );
}
