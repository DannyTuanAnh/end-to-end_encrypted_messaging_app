export default function Footer() {
  return (
    <footer className="border-t h-16 px-4 py-2 text-center text-sm text-muted-foreground">
      © {new Date().getFullYear()} Chat App — built with React + Vite
    </footer>
  );
}
