"use client";

import Link from "next/link";
import { useUser, SignedIn, SignedOut, SignOutButton } from "@clerk/nextjs";
import { Button } from "@/components/ui/button"; // prilagodjeno ako nema src foldera

export default function HomePage() {
  const { user } = useUser();

  return (
    <div className="min-h-screen bg-gray-50">
      <main className="p-8">
        <h1 className="text-3xl font-bold mb-4">
          Welcome to Document Intelligence
        </h1>
        <p className="mb-6">
          This is your AI-powered document management and chat platform.
        </p>

        <SignedIn>
          <div className="mt-6 p-4 border rounded bg-white shadow">
            <p>Hello, {user?.firstName || user?.username}!</p>
            <p>You are logged in and ready to manage your documents.</p>
          </div>
        </SignedIn>

        <SignedOut>
          <div className="mt-6 p-4 border rounded bg-white shadow">
            <p>Please log in to access upload and chat features.</p>
          </div>
        </SignedOut>
      </main>
    </div>
  );
}
