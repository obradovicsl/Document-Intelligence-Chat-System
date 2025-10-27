"use client";
import Link from "next/link";
import { useUser, SignedIn, SignedOut, SignOutButton } from "@clerk/nextjs";
import { Button } from "@/components/ui/button";

export default function Navbar() {
  const { user } = useUser();

  return (
    <nav className="flex justify-between items-center p-4 bg-white shadow">
      <Link href="/">
        <div className="text-xl font-bold">DocIntelligence</div>
      </Link>
      <div className="space-x-4">
        <SignedOut>
          <Link href="/sign-in">
            <Button variant="outline">Login</Button>
          </Link>
          <Link href="/sign-up">
            <Button variant="outline">Sign Up</Button>
          </Link>
        </SignedOut>
        <SignedIn>
          <span className="mr-2">{user?.firstName || user?.username}</span>
          <SignOutButton>
            <Button variant="destructive">Logout</Button>
          </SignOutButton>
        </SignedIn>
      </div>
    </nav>
  );
}
