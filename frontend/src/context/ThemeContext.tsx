import { createContext, useContext, useState } from "react";

export type Theme = "light" | "dark";

export type ThemeContextType = {
  theme: Theme;
  toggleTheme: () => void;
};

const ThemeContext = createContext<ThemeContextType>({
  theme: "light",
  toggleTheme: () => {},
});

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [theme, setTheme] = useState<Theme>(() => {
    try {
      const storedTheme = localStorage.getItem("theme") as Theme | null;
      return storedTheme || "light";
    } catch (error) {
      console.error("Failed to retrieve theme from localStorage", error);
      return "light";
    }
  });

  const toggleTheme = () => {
    setTheme((prevTheme) => (prevTheme === "light" ? "dark" : "light"));
    try {
      localStorage.setItem("theme", theme === "light" ? "dark" : "light");
    } catch (error) {
      console.error("Failed to save theme to localStorage", error);
    }
  };

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => useContext(ThemeContext);
