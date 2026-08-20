"use client";

import { useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

interface IProps {
  email: string;
  onBack: () => void;
}

const OtpVerification = ({ email, onBack }: IProps) => {
  const [otp, setOtp] = useState(["", "", "", "", "", ""]);
  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);
  const router = useRouter();

  const handleChange = (value: string, index: number) => {
    if (!/^\d*$/.test(value)) return;

    const newOtp = [...otp];
    newOtp[index] = value.slice(-1);
    setOtp(newOtp);

    if (value && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handleKeyDown = (
    e: React.KeyboardEvent<HTMLInputElement>,
    index: number,
  ) => {
    if (e.key === "Backspace" && !otp[index] && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
  };

  const handlePaste = (e: React.ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault();

    const pastedOtp = e.clipboardData
      .getData("text")
      .replace(/\D/g, "")
      .slice(0, 6);

    if (!pastedOtp) return;

    const newOtp = ["", "", "", "", "", ""];

    pastedOtp.split("").forEach((digit, index) => {
      newOtp[index] = digit;
    });

    setOtp(newOtp);

    const nextIndex = Math.min(pastedOtp.length, 5);
    inputRefs.current[nextIndex]?.focus();
  };

  const handleLogin = async () => {
    const otpValue = otp.join("");

    if (otpValue.length !== 6) {
      toast.error("Please enter the complete 6-digit OTP");
      return;
    }

    try {
      const response = await fetch("/api/auth/login", {
        method: "POST",
        body: JSON.stringify({
          otp: otpValue,
          email,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      });

      const data = await response.json();

      if (!response.ok) {
        toast.error(data.message);
        return;
      }

      toast.success(data.message);
      sessionStorage.removeItem("verify_email");
      router.push("/");
    } catch (error) {
      console.error(error);
      toast.error("Internal server error");
    }
  };

  return (
    <div className="w-full">
      {/* Heading */}
      <div>
        <h1 className="font-serif text-[38px] font-bold leading-tight text-[#17191c]">
          Enter the OTP
        </h1>

        <p className="mt-3 text-[16px] leading-7 text-[#686c73]">
          Enter the 6-digit code sent to your email address
          <br />
          to verify your account.
        </p>
      </div>

      {/* OTP Inputs */}
      <div className="mt-9">
        <label className="mb-4 block text-[15px] font-medium text-[#292b2f]">
          OTP Code
        </label>

        <div className="flex gap-3">
          {otp.map((digit, index) => (
            <input
              key={index}
              ref={(el) => {
                inputRefs.current[index] = el;
              }}
              type="text"
              inputMode="numeric"
              maxLength={1}
              value={digit}
              onChange={(e) => handleChange(e.target.value, index)}
              onKeyDown={(e) => handleKeyDown(e, index)}
              onPaste={handlePaste}
              className="
                h-[60px]
                w-[58px]
                rounded-xl
                border
                border-[#dedede]
                bg-white
                text-center
                text-[22px]
                font-medium
                text-[#17191c]
                outline-none
                transition-all

                focus:border-[#b45f32]
                focus:ring-2
                focus:ring-[#b45f32]/15
              "
            />
          ))}
        </div>
      </div>

      {/* Expiration */}
      <div className="mt-7 flex items-center gap-2 text-[14px] text-[#73777d]">
        <svg
          width="18"
          height="18"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.8"
        >
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z" />
          <path d="M12 8v4l2.5 2.5" />
        </svg>

        <span>
          The OTP will expire in{" "}
          <span className="font-semibold text-[#b45f32]">05:00</span>
        </span>
      </div>

      {/* Verify Button */}
      <button
        onClick={handleLogin}
        className="
          mt-8
          h-[58px]
          w-full
          rounded-xl
          bg-[#b45f32]
          text-[17px]
          font-semibold
          text-white
          transition-all

          hover:bg-[#9f4f2c]
          active:scale-[0.99]
        "
      >
        Verify OTP
      </button>

      {/* Divider */}
      <div className="my-7 flex items-center gap-5">
        <div className="h-px flex-1 bg-[#e5e5e5]" />

        <span className="text-[14px] text-[#777b81]">or</span>

        <div className="h-px flex-1 bg-[#e5e5e5]" />
      </div>

      {/* Back to Login */}
      <button
        onClick={onBack}
        className="
          flex
          h-[55px]
          w-full
          items-center
          justify-center
          gap-2
          rounded-xl
          border
          border-[#dedede]
          bg-white
          text-[16px]
          font-medium
          text-[#303238]
          transition-all

          hover:bg-[#fafafa]
        "
      >
        <span className="text-xl">←</span>
        Back to login
      </button>
    </div>
  );
};

export default OtpVerification;
