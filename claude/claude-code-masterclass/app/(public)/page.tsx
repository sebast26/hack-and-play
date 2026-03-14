import Link from "next/link";
import { Clock8, Target, Users, Zap, ArrowRight } from "lucide-react";
import styles from "./page.module.css";

export default function Home() {
  return (
    <div className={styles.splash}>
      <div className={styles.grid} aria-hidden="true" />
      <div className={styles.glow} aria-hidden="true" />

      <div className={styles.hero}>
        <div className={styles.badge}>
          <span className={styles.dot} />
          Mission Status: Active
        </div>

        <h1 className={styles.title}>
          P<Clock8 className={styles.clock} strokeWidth={2.5} />
          cket
          <span className={styles.heist}>Heist</span>
        </h1>

        <p className={styles.tagline}>
          Mischief on the clock.
          <span className={styles.cursor} aria-hidden="true" />
        </p>

        <p className={styles.description}>
          Welcome to Pocket Heist — the app that turns your office into a
          playground. Plan covert missions, assign tasks to unsuspecting
          colleagues, and track the chaos as it unfolds. No heist too small, no
          mischief too grand.
        </p>

        <ul className={styles.features}>
          <li className={styles.featureItem}>
            <Target size={14} />
            Plan covert missions
          </li>
          <li className={styles.featureItem}>
            <Users size={14} />
            Assign unsuspecting colleagues
          </li>
          <li className={styles.featureItem}>
            <Zap size={14} />
            Track the chaos live
          </li>
        </ul>

        <div className={styles.actions}>
          <Link href="/signup" className={styles.cta}>
            Enlist Now
            <ArrowRight size={16} />
          </Link>
          <Link href="/login" className={styles.loginLink}>
            Already operative? <span>Log in</span>
          </Link>
        </div>
      </div>

      <div className={styles.scanline} aria-hidden="true" />
    </div>
  );
}
