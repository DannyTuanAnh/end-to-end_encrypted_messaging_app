import Header from "./Header";
import Footer from "./Footer";
import { Outlet } from "react-router-dom";

export default function MainLayout() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 p-4 bg-background overflow-auto flex justify-center items-center">
        <Outlet />
      </main>
      <Footer />
    </div>
  );
}
