## Shared Components

This file lists the shared components found in `src/components/` and explains their purpose and common props.

Components live under `src/components/` and are organized into:

- `src/components/ui/` — shadcn primitives and small wrappers
- `src/components/common/` — reusable presentational components (AvatarUser, ChatItem, UserCard)
- `src/components/chat/` — chat-specific components (MessageBubble, ChatInput, ChatHeader)

- `Avatar` (avatar.tsx)
  - Purpose: Render user's avatar with fallback initials.
  - Props:
    - `src?: string` — image url
    - `alt?: string` — alt text
    - `size?: 'sm' | 'md' | 'lg'` — size variant
  - Usage:

```tsx
<Avatar src={user.avatar} alt={user.name} size="md" />
```

- `Badge` (common/badge.tsx)
  - Purpose: Small status or count pill.
  - Props: `children`, `variant?` (color)

- `Button` (ui/button.tsx)
  - Purpose: Styled button wrapper around native `button`.
  - Props: `onClick`, `children`, `variant?`, `disabled?`.
  - Usage:

```tsx
<Button variant="primary" onClick={submit}>
  Send
</Button>
```

- `Card` (common/card.tsx)
  - Purpose: Elevated container used for grouped content.
  - Props: `children`, `className`.

- `Input` / `Textarea` (ui/input.tsx, ui/textarea.tsx)
  - Purpose: Form controls with consistent styling and labels.
  - Props: `value`, `onChange`, `placeholder`, `name`, `error?`.

- `Dialog` (ui/dialog.tsx)
  - Purpose: Modal dialog using shadcn primitives.
  - Props: `open`, `onOpenChange`, `title`, `children`.

- `DropdownMenu` / `ContextMenu` (ui/dropdown-menu.tsx, ui/context-menu.tsx)
  - Purpose: Action menus and contextual commands.

- `Switch`, `Tabs`, `Tooltip`, `Separator`, `Spinner` (ui/switch.tsx, ui/tabs.tsx, ui/tooltip.tsx, ui/separator.tsx, ui/skeleton.tsx)
  - Purpose: UI primitives for interactivity and layout.

Guidelines for using shared components

- Pass only serializable props to components when possible.
- Keep layout-specific props like `className` available for overrides.
- Prefer composition over prop bloat: pass children or render props for complex content.

When to add a new shared component

- When UI is reused in 2+ places.
- When a design needs consistent behavior (loading states, a11y, focus handling).
