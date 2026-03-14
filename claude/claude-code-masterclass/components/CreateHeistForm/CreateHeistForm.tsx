"use client";

import { useState, useEffect, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import {
  collection,
  getDocs,
  addDoc,
  serverTimestamp,
} from "firebase/firestore";
import { db } from "@/lib/firebase/config";
import { useUser } from "@/context/AuthContext";
import {
  COLLECTIONS,
  type UserDoc,
  type CreateHeistInput,
} from "@/types/firestore";
import styles from "./CreateHeistForm.module.css";

export default function CreateHeistForm() {
  const { user } = useUser();
  const router = useRouter();

  const [users, setUsers] = useState<UserDoc[]>([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [selectedUser, setSelectedUser] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchUsers() {
      const snapshot = await getDocs(collection(db, COLLECTIONS.USERS));
      const docs = snapshot.docs.map((doc) => doc.data() as UserDoc);
      setUsers(docs);
    }
    fetchUsers();
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      const [assignedTo, assignedToCodename] = selectedUser.split("|");
      const deadline = new Date(Date.now() + 48 * 60 * 60 * 1000);

      const payload: CreateHeistInput = {
        title,
        description,
        createdBy: user!.uid,
        createdByCodename: user!.displayName ?? "",
        assignedTo,
        assignedToCodename,
        createdAt: serverTimestamp(),
        deadline,
        finalStatus: null,
      };

      await addDoc(collection(db, COLLECTIONS.HEISTS), payload);
      router.push("/heists");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Something went wrong";
      setError(message);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Mission title"
        required
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className={styles.input}
      />

      <textarea
        placeholder="Describe the mission..."
        required
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        className={styles.textarea}
        rows={4}
      />

      <select
        required
        value={selectedUser}
        onChange={(e) => setSelectedUser(e.target.value)}
        className={styles.select}
      >
        <option value="">
          {users.length === 0
            ? "No operatives available"
            : "Assign to operative..."}
        </option>
        {users.map((u) => (
          <option key={u.id} value={`${u.id}|${u.codename}`}>
            {u.codename}
          </option>
        ))}
      </select>

      {error && <p className={styles.error}>{error}</p>}

      <button
        type="submit"
        className="btn"
        disabled={isLoading || users.length === 0}
      >
        {isLoading ? "Launching..." : "Launch Heist"}
      </button>
    </form>
  );
}
