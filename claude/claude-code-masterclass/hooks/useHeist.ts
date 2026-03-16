"use client";

import { useState, useEffect } from "react";
import { collection, query, where, onSnapshot } from "firebase/firestore";
import { db } from "@/lib/firebase/config";
import { useUser } from "@/context/AuthContext";
import { COLLECTIONS, heistConverter, type Heist } from "@/types/firestore";

export type HeistMode = "active" | "assigned" | "expired";

interface UseHeistResult {
  heists: Heist[];
  loading: boolean;
  error: string | null;
}

export function useHeist(mode: HeistMode): UseHeistResult {
  const { user } = useUser();
  const [heists, setHeists] = useState<Heist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!user) {
      setLoading(false);
      return;
    }

    const now = new Date();
    const col = collection(db, COLLECTIONS.HEISTS).withConverter(
      heistConverter,
    );

    let q;
    if (mode === "active") {
      q = query(
        col,
        where("assignedTo", "==", user.uid),
        where("deadline", ">", now),
      );
    } else if (mode === "assigned") {
      q = query(
        col,
        where("createdBy", "==", user.uid),
        where("deadline", ">", now),
      );
    } else {
      q = query(col, where("deadline", "<=", now));
    }

    const unsubscribe = onSnapshot(
      q,
      (snapshot) => {
        const data = snapshot.docs.map((doc) => doc.data()) as Heist[];
        const filtered =
          mode === "expired"
            ? data.filter((h) => h.finalStatus !== null)
            : data;
        setHeists(filtered);
        setLoading(false);
      },
      (err) => {
        setError(err.message);
        setLoading(false);
      },
    );

    return unsubscribe;
  }, [user, mode]);

  return { heists, loading, error };
}
