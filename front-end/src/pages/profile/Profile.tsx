import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function Profile() {
  return (
    <div className="min-w-sm w-full max-w-3xl">
      <h1 className="text-2xl font-semibold mb-4">Profile</h1>
      <Card className="p-4">
        <div className="flex items-center gap-4">
          <Avatar>
            <AvatarImage src="/assets/avatar-placeholder.png" alt="User" />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
          <div>
            <div className="text-lg font-medium">User Name</div>
            <div className="text-sm text-muted-foreground">
              user@example.com
            </div>
          </div>
        </div>

        <div className="mt-4">
          <Button>Update profile</Button>
        </div>
      </Card>
    </div>
  );
}
