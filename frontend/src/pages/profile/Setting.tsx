import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { useTheme } from "@/context/ThemeContext";
export default function Setting() {
  const { theme, toggleTheme } = useTheme();
  return (
    <Card className="p-4">
      <CardHeader>
        <CardTitle>Settings</CardTitle>
        <CardDescription>
          This section can include user-specific settings and preferences.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <p>Change your settings here.</p>
        <Separator />
        <div className="flex items-center justify-between">
          <p>Change theme</p>
          <div className="flex items-center space-x-2">
            <Switch
              id="airplane-mode"
              checked={theme === "dark"}
              onCheckedChange={toggleTheme}
            />
            <Label htmlFor="airplane-mode">{theme}</Label>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
