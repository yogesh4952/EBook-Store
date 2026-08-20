"use client";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

const page = () => {
  const [otp, setOtp] = useState("");
  const router = useRouter();
  const email = sessionStorage.getItem("verify_email");
  const handleLogin = async () => {
    try {
      const payload = {
        otp,
        email,
      };
      const response = await fetch("/api/auth/login", {
        method: "POST",
        body: JSON.stringify(payload),
        headers: {
          "Content-type": "application/json",
        },
      });
      const data = await response.json();
      if (!response.ok) {
        toast.error(data.message);
      }
      router.push("/");
      toast.success(data.message);
      sessionStorage.removeItem("verify_email");
    } catch (error) {
      toast.error("Internal server error");
      console.error(error);
    }
  };
  return (
    <div>
      <div>
        <label htmlFor="otp">Enter OTP</label>
        <input
          type="text"
          name="otp"
          id="otp"
          onChange={(e) => setOtp(e.target.value)}
          value={otp}
        />
      </div>

      <button onClick={handleLogin}>Login</button>
    </div>
  );
};

export default page;
