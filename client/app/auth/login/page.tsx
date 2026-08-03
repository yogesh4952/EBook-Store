import { FaFacebook, FaGoogle } from "react-icons/fa";
import { FaShield } from "react-icons/fa6";

const Login = () => {
  return (
    <div className="grid grid-cols-2 pt-8 gap-6">
      <div
        className="w-full h-full rounded-xl bg-cover bg-center bg-no-repeat"
        style={{ backgroundImage: "url('/login.jpg')" }}
      ></div>
      <div className="">
        <h1 className="text-3xl font-semibold">Login to your account</h1>
        <p className="text-gray-500 mt-4">
          Enter your credentials to access your account.
        </p>

        {/* form */}
        <div className="flex flex-col gap-8 mt-8">
          <div className="flex flex-col">
            <label className="" htmlFor="email">
              Email:
            </label>
            <input
              className="w-full border border-accent rounded px-4 py-2"
              type="text"
              name="email"
              id="email"
            />
          </div>

          <div className="flex flex-col">
            <label htmlFor="email">Password:</label>
            <input
              className="w-full border border-accent rounded px-4 py-2"
              type="text"
              name="email"
              id="email"
            />
          </div>
        </div>

        <p className="text-primary text-right mt-4 font-bold cursor-pointer">
          Forgot Password?
        </p>

        <div>
          <button className="w-full rounded bg-primary text-white h-12 mt-8 cursor-pointer hover:bg-amber-950">
            Login
          </button>

          <p className="text-gray-600 text-center mt-10">or continue with</p>

          {/* Oauth */}
          <div className="mt-10 gap-6 grid">
            <button className="flex gap-4 items-center text-center justify-center w-full rounded h-12 bg-white shadow border border-gray-300">
              <FaGoogle />
              Continue with Google
            </button>

            <button className="flex gap-4 items-center text-center justify-center w-full rounded h-12 bg-white shadow border border-gray-300">
              <FaFacebook />
              Continue with Facebook
            </button>
          </div>

          <p className="text-gray-400 flex items-center text-center mt-10 gap-2">
            <FaShield />
            We never share your data with anyone
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
