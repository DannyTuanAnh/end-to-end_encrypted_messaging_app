export type Room = {
  id: string;
  name: string;
  avatar: string;
  participants: string[]; // userIds
  lastMessage: string;
  lastMessageAt: string;
};

export const dataRoom: Room[] = [
  {
    id: "room1",
    name: "Alice",
    avatar: "/avatars/alice.png",
    participants: ["u1", "u2"],
    lastMessage: "I'm good, thanks!",
    lastMessageAt: "2026-04-01T10:00:00Z",
  },
  {
    id: "room2",
    name: "Support",
    avatar: "/avatars/support.png",
    participants: ["u1", "support"],
    lastMessage: "Can you provide your order ID?",
    lastMessageAt: "2026-04-01T09:30:00Z",
  },
  {
    id: "room3",
    name: "Charlie",
    avatar: "/avatars/charlie.png",
    participants: ["u1", "u3"],
    lastMessage: "That last-minute goal was crazy!",
    lastMessageAt: "2026-04-01T08:45:00Z",
  },
];
