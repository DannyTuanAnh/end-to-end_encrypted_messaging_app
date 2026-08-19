import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

import MainLayout from "@/layouts/MainLayout";
import AuthLayout from "@/layouts/AuthLayout";
import Landing from "@/pages/Landing";
import Home from "@/pages/Home";
import Profile from "@/pages/profile/Profile";
import ChatLayout from "@/layouts/ChatLayout";
import ChatList from "@/pages/chat/ChatList";
import ChatPage from "@/pages/chat/ChatPage";
import Login from "@/pages/auth/Login";

function RequireAuth({ children }: { children: React.ReactElement }) {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  if (!token) {
    return <Navigate to="/auth/login" replace />;
  }
  return children;
}

export default function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/auth" element={<AuthLayout />}>
          <Route index element={<Navigate to="login" replace />} />
          <Route path="login" element={<Login />} />
        </Route>
        <Route path="/" element={<Landing />} />
        <Route
          element={
            <RequireAuth>
              <MainLayout />
            </RequireAuth>
          }
        >
          {/* <Route index element={<Home />} /> */}
          <Route path="home" element={<Home />} />
          <Route path="profile" element={<Profile />} />
        </Route>
        {/* Chat routes: list and selected chat using ChatLayout (4:6 split) */}
        <Route path="chat" element={<ChatLayout />}>
          <Route index element={<ChatList />} />
          <Route path=":id" element={<ChatPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
