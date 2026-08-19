import { Button } from "@/components/ui/button";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { useGoogleLogin } from "@react-oauth/google";

import { Home } from "lucide-react";

import { Link, useNavigate } from "react-router-dom";

import { useAuthContext } from "@/context/AuthContext";

import { useState } from "react";

export default function Login() {
  const navigate = useNavigate();

  const { login } = useAuthContext();

  const [loading, setLoading] = useState(false);

  const googleLogin = useGoogleLogin({
    flow: "auth-code",

    onSuccess: async (codeResponse) => {
      try {
        setLoading(true);

        /**
         * codeResponse.code
         * chính là auth_code backend yêu cầu
         */
        await login(codeResponse.code);

        navigate("/home");
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    },

    onError: () => {
      console.error("Google Login Failed");
    },
  });

  async function handleGoogleLogin() {
    googleLogin();
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Login</CardTitle>

        <CardDescription>
          Welcome back! Please login to your account.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <Button
          onClick={handleGoogleLogin}
          disabled={loading}
          className="w-full"
        >
          {loading ? "Logging in..." : "Login with Google"}
        </Button>
      </CardContent>

      <CardFooter>
        <Link to="/" className="w-full">
          <Button variant="ghost" className="w-full">
            <Home />
            Back to Home
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
