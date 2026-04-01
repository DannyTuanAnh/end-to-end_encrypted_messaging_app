export type Message = {
  id: string;
  roomId: string;
  senderId: string;
  content: string;
  createdAt: string;
};

export const dataMessage: Record<string, Message[]> = {
  room1: [
    {
      id: "m1",
      roomId: "room1",
      senderId: "u2",
      content: "Hi there!",
      createdAt: "2026-04-01T09:58:00Z",
    },
    {
      id: "m2",
      roomId: "room1",
      senderId: "u1",
      content: "Hello! How are you?",
      createdAt: "2026-04-01T09:59:00Z",
    },
    {
      id: "m3",
      roomId: "room1",
      senderId: "u2",
      content: "I'm good, thanks!",
      createdAt: "2026-04-01T10:00:00Z",
    },
  ],

  room2: [
    {
      id: "m4",
      roomId: "room2",
      senderId: "support",
      content: "Welcome to support chat!",
      createdAt: "2026-04-01T09:20:00Z",
    },
    {
      id: "m5",
      roomId: "room2",
      senderId: "u1",
      content: "I need help with my order.",
      createdAt: "2026-04-01T09:25:00Z",
    },
    {
      id: "m6",
      roomId: "room2",
      senderId: "support",
      content: "Can you provide your order ID?",
      createdAt: "2026-04-01T09:30:00Z",
    },
  ],

  room3: [
    {
      id: "m7",
      roomId: "room3",
      senderId: "u3",
      content: "Did you see the game last night?",
      createdAt: "2026-04-01T08:40:00Z",
    },
    {
      id: "m8",
      roomId: "room3",
      senderId: "u1",
      content: "Yeah, it was amazing!",
      createdAt: "2026-04-01T08:42:00Z",
    },
    {
      id: "m9",
      roomId: "room3",
      senderId: "u3",
      content: "That last-minute goal was crazy!",
      createdAt: "2026-04-01T08:45:00Z",
    },
  ],
};
