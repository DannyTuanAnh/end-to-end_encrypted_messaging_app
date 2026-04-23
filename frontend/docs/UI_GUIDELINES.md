## UI Guidelines

This section describes the main UI layout, components used from shadcn/ui and styling conventions.

UI structure

- Sidebar: navigation and channel list. Collapsible, contains avatars and channel badges.
- Header: top bar with current chat name, actions and user menu.
- Chat area: message list (scrollable) and message input footer.
- Message list: grouped by day, newest messages at bottom, virtualized when large.

shadcn/ui primitives used

- `Dialog` — modal flows (profile, settings).
- `DropdownMenu` — user actions.
- `Input` / `Textarea` — message compose and forms.
- `Tabs` / `Switch` / `Tooltip` — miscellaneous controls.

Styling conventions

- Use Tailwind utility classes for layout and spacing.
- Keep component-specific classes inside component files. Expose `className` prop for overrides.
- Prefer semantic HTML and accessible attributes (`aria-*`, keyboard handlers).

Responsive rules

- Mobile: Sidebar collapses into a drawer; chat list accessible via icon.
- Medium / Desktop: persistent sidebar visible; two-panel layout: channels and chat area.

Accessibility

- All actionable elements receive `aria-label` or descriptive text.
- Ensure color contrast for primary text and interactive elements.

Design tokens

- Keep colors, spacing and font sizes in Tailwind config where needed. Prefer utility classes for rapid iteration.
