## Folder Structure — Frontend

This document explains the main folders in the frontend and their responsibilities.

```
src/
├── assets/              # images, icons, fonts
├── components/          # shared components
│   ├── ui/              # shadcn components
│   ├── common/          # reusable components (AvatarUser, ChatItem...)
│   └── chat/            # chat specific components (Message, ChatInput...)
│
├── pages/               # route pages
│   ├── auth/
│   │   ├── Login.tsx
│   │   └── Register.tsx
│   ├── chat/
│   │   └── ChatPage.tsx
│   └── profile/
│       └── Profile.tsx
│
├── layouts/             # app layouts
│   ├── MainLayout.tsx
│   └── AuthLayout.tsx
│
├── hooks/               # custom hooks
│   ├── useSocket.ts
│   ├── useAuth.ts
│   └── useDebounce.ts
│
├── services/            # API calls / axios
│   ├── axios.ts
│   ├── auth.service.ts
│   ├── user.service.ts
│   └── message.service.ts
│
├── store/               # Zustand stores
│   ├── authStore.ts
│   ├── chatStore.ts
│   └── uiStore.ts
│
├── socket/              # socket.io logic
│   └── socket.ts
│
├── lib/                 # utilities
│   ├── utils.ts
│   ├── constants.ts
│   └── format.ts
│
├── types/               # TypeScript types
│   ├── user.ts
│   ├── message.ts
│   └── api.ts
│
├── routes/              # router config
│   └── router.tsx
│
├── styles/              # global css
│   └── index.css
│
├── App.tsx
└── main.tsx
```

**src/components/**

Purpose: Shared UI components and reusable presentational components.
Structure:

- ui/ → shadcn UI components
- common/ → reusable components (AvatarUser, ChatItem, UserCard)
- chat/ → chat related components (MessageBubble, ChatInput, ChatHeader)

Rule: Components should be presentational and receive data via props.

**src/pages/**

Purpose: Route-level pages. Each page represents a route in the application.

Examples:

- auth/Login.tsx
- auth/Register.tsx
- chat/ChatPage.tsx
- profile/Profile.tsx

Pages compose components, hooks, stores and services.

**src/layouts/**

Purpose: Layout wrappers used by routes.

Examples:

- MainLayout → sidebar + chat layout
- AuthLayout → login/register layout
- src/hooks/

Purpose: Reusable custom hooks that contain logic and side effects.

Examples:

- useSocket
- useAuth
- useDebounce
- useMobile

Keep logic here instead of inside components.

**src/services/**

Purpose: API service layer and HTTP logic.

Contains:

- Axios instance configuration
- API endpoint functions
- Request/response handling
- Token handling

Example:

- auth.service.ts
- message.service.ts
- user.service.ts

**src/store/**

Purpose: Global client state using Zustand.

example:

- authStore → user, token
- chatStore → messages, rooms, current chat
- uiStore → sidebar, theme, modal
  **src/socket/**

Purpose: Socket.io connection and realtime event handling.

Examples:

- connect socket
- send message
- receive message
- typing indicator
- online users

**src/lib/ or src/utils/**

Purpose: Utility functions and helpers.

Examples:

- format date
- debounce
- validation
- constants
- helper functions

**src/types/**

Purpose: Global TypeScript types and interfaces.

Examples:

- User
- Message
- ChatRoom
- API Response

**src/routes/**

Purpose: Routing configuration and protected routes.

Contains:

- Route definitions
- Layout wrapping
- Auth protected routes

**src/assets/**

Purpose: Static assets such as images, icons and fonts.

**src/styles/**

Purpose: Global styles and Tailwind entry files.
