import { createContext, useContext, useEffect, useState } from "react";

import { getProfile, loginGoogle, logoutApi, type User } from "@/api/auth.api";

type AuthContextType = {
  user: User | null;
  loading: boolean;
  isAuthenticated: boolean;

  login: (authCode: string) => Promise<void>;
  logout: () => Promise<void>;

  refreshProfile: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);

  const [loading, setLoading] = useState(true);

  async function refreshProfile() {
    try {
      const profile = await getProfile();

      setUser(profile);
    } catch {
      setUser(null);
    }
  }

  async function login(authCode: string) {
    await loginGoogle({
      auth_code: authCode,
    });

    await refreshProfile();
  }

  async function logout() {
    try {
      await logoutApi();
    } finally {
      setUser(null);
    }
  }

  useEffect(() => {
    async function initAuth() {
      try {
        await refreshProfile();
      } finally {
        setLoading(false);
      }
    }

    initAuth();
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        isAuthenticated: !!user,

        login,
        logout,

        refreshProfile,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuthContext must be used within AuthProvider");
  }

  return context;
}
