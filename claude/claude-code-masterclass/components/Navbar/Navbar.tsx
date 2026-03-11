"use client";

import { Clock8, LogOut } from "lucide-react";
import Link from "next/link";
import { signOut } from "firebase/auth";
import Avatar from "@/components/Avatar";
import { useUser } from "@/context/AuthContext";
import { auth } from "@/lib/firebase/config";
import styles from "./Navbar.module.css";

export default function Navbar() {
  const { user, isLoading } = useUser();

  async function handleLogout() {
    try {
      await signOut(auth);
    } catch (err) {
      console.error("Logout failed", err);
    }
  }

  return (
    <div className={styles.siteNav}>
      <nav>
        <header>
          <h1>
            <Link href="/heists">
              P<Clock8 className={styles.logo} size={14} strokeWidth={2.75} />
              cket Heist
            </Link>
          </h1>
          <div>Tiny missions. Big office mischief.</div>
        </header>
        <ul>
          {!isLoading && user && (
            <li>
              <button className={styles.logoutBtn} onClick={handleLogout}>
                <LogOut size={16} />
                Logout
              </button>
            </li>
          )}
          <li>
            <Link href="/heists/create" className="btn">
              Create Heist
            </Link>
          </li>
          {!isLoading && user && (
            <li>
              <Avatar name={user.displayName ?? user.email ?? "User"} />
            </li>
          )}
        </ul>
      </nav>
    </div>
  );
}
