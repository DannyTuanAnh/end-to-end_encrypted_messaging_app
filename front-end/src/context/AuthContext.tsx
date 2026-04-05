import React, { createContext, useContext, useEffect, useState } from "react";

export type User = {
  email?: string;
  name?: string;
};

export type AuthReturn = {
  token: string | null;
  user: User | null;
  login: (email: string, password?: string) => void;
  logout: () => void;
  isAuthenticated: boolean;
};

// Demo auth hook: keeps state in React and mirrors to localStorage.
export function useAuth(): AuthReturn {
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
      return raw ? JSON.parse(raw) : null;
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
    const demoToken = "demo-token";
    const u: User = { email };
    try {
      localStorage.setItem("token", demoToken);
      localStorage.setItem("user", JSON.stringify(u));
    } catch {}
    setToken(demoToken);
    setUser(u);
  };

  const logout = () => {
    try {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    } catch {}
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

// Context and provider
export const AuthContext = createContext<AuthReturn | undefined>(undefined);

export function useAuthContext(): AuthReturn {
  const ctx = useContext(AuthContext);
  if (!ctx)
    throw new Error("useAuthContext must be used within an AuthProvider");
  return ctx;
}

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const auth = useAuth();
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
};

export default AuthProvider;
