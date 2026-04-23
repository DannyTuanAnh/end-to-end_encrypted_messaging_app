import React, { createContext, useContext, useEffect, useState } from "react";

export type User = {
  email?: string;
  name?: string;
  uid: string;
};

export type AuthReturn = {
  token: string | null;
  user: User | null;
  login: (email: string, password?: string) => void;
  logout: () => void;
  isAuthenticated: boolean;
};

const AuthContext = createContext<AuthReturn | undefined>(undefined);

// Demo auth hook: keeps state in React and mirrors to localStorage.
function useAuth(): AuthReturn {
  const [token, setToken] = useState<string | null>(() => {
    try {
      return localStorage.getItem("token");
    } catch {
      return null;
    }
  });

  const [user, setUser] = useState<User | null>(() => {
    try {
      const raw = localStorage.getItem("user");
      const u = raw ? JSON.parse(raw) : null;
      if (u.uid === "") {
        return console.log(
          "Invalid user data in localStorage: uid is empty string",
        );
      }
      return u;
    } catch {
      return null;
    }
  });

  useEffect(() => {
    function onStorage(e: StorageEvent) {
      if (e.key === "token") setToken(e.newValue);
      if (e.key === "user") setUser(e.newValue ? JSON.parse(e.newValue) : null);
    }

    window.addEventListener("storage", onStorage);
    return () => window.removeEventListener("storage", onStorage);
  }, []);

  const login = (email: string, _password?: string) => {
    const demoToken = email && _password ? `token-for-${email}` : null;
    const u: User = {
      email,
      name: email.split("@")[0],
      uid: `u_${email.split("@")[0]}`,
    };
    if (!demoToken) return;
    try {
      localStorage.setItem("token", demoToken);
      localStorage.setItem("user", JSON.stringify(u));
    } catch (error) {
      console.error("Failed to save auth data to localStorage", error);
    }
    setToken(demoToken);
    setUser(u);
  };

  const logout = () => {
    try {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    } catch (error) {
      console.error("Failed to remove auth data from localStorage", error);
    }
    setToken(null);
    setUser(null);
  };

  return {
    token,
    user,
    login,
    logout,
    isAuthenticated: !!token,
  };
}

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const auth = useAuth();
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
};

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx)
    throw new Error("useAuthContext must be used within an AuthProvider");
  return ctx;
}
