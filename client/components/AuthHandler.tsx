"use client";

import { useSearchParams } from "next/navigation";
import { useEffect } from "react";
import { toast } from "sonner";

const AuthHandler = () => {
  const searchParams = useSearchParams();

  useEffect(() => {
    const err = searchParams.get("error");
    if (err === "AuthenticationFailed") {
      toast.error("Google login failed due to server error (duplicate phone / DB issue).");
    } else if (err === "DBConflict") {
      toast.error("Account creation failed: phone number conflict.");
    } else if (err) {
      toast.error("Login error: " + err);
    }
  }, [searchParams]);

  return null;
};

export default AuthHandler;
