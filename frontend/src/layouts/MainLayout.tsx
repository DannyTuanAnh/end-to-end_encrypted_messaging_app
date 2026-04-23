import Header from "./Header";
import Footer from "./Footer";
import { Outlet } from "react-router-dom";

export default function MainLayout() {
  return (
    <div className="h-screen flex flex-col">
      <Header />
      <main className="mt-16 p-8 bg-background overflow-auto flex justify-center items-start h-[calc(100vh-4rem)]">
        <Outlet />
      </main>
      <Footer />
    </div>
  );
}
