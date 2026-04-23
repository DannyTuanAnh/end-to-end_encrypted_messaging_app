import ChatList from "../pages/chat/ChatList";
import ChatPage from "../pages/chat/ChatPage";
import Header from "@/layouts/Header";

// ChatLayout provides a 4:6 two-column layout: left ChatList, right Outlet (ChatPage)
export default function ChatLayout() {
  return (
    <div className="flex w-full overflow-hidden">
      <Header></Header>
      <div className="w-2/5 min-w-[260px] min-h-screen border-r">
        <ChatList />
      </div>

      <div className="w-3/5 h-screen p-6 overflow-hidden">
        <div className="h-full p-0">
          <ChatPage />
        </div>
      </div>
    </div>
  );
}
