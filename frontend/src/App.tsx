import Router from "@/routes/router";
import { AuthProvider } from "@/context/AuthContext";
import { ThemeProvider } from "@/context/ThemeContext";
import axios from "axios";
import { useEffect } from "react";
const checkConnection = async () => {
  try {
    // Endpoint này mở công khai và không cần API Key hay Auth
    const response = await axios.get("https://chat-app-ta.duckdns.org/healthz");
    console.log("✅ Kết nối Backend thành công! Trạng thái:", response.data);
  } catch (error) {
    console.error("❌ Kết nối Backend thất bại!", error);
  }
};
function App() {
  useEffect(() => {
    checkConnection();
  }, []);
  return (
    <AuthProvider>
      <ThemeProvider>
        <Router />
      </ThemeProvider>
    </AuthProvider>
  );
}

export default App;
