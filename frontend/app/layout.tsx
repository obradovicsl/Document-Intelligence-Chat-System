import { ClerkProvider } from "@clerk/nextjs";
import "./globals.css";
import Navbar from "../components/Navbar";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ClerkProvider
      publishableKey={process.env.NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY!}
    >
      <html lang="en">
        <body className="min-h-screen">
          <Navbar />

          {children}
        </body>
      </html>
    </ClerkProvider>
  );
}
