"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Clock8 } from "lucide-react";
import { useUser } from "@/context/AuthContext";

export default function PublicLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const { user, isLoading } = useUser();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && user) {
      router.replace("/heists");
    }
  }, [isLoading, user, router]);

  if (isLoading || user) {
    return (
      <div className="auth-loader">
        <Clock8 size={32} strokeWidth={2.75} />
      </div>
    );
  }

  return <main className="public">{children}</main>;
}
