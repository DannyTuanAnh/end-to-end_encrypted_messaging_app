import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardHeader,
  CardContent,
  CardTitle,
  CardDescription,
  CardAction,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/context/AuthContext";
import { DialogTrigger, Dialog } from "@/components/ui/dialog";
import { EditProfile } from "@/components/profile/EditProfile";
import { Edit2 } from "lucide-react";
import QRCode from "@/components/common/QRCode";

export default function Profile() {
  const { user } = useAuth();
  const uid = "user_123456";
  return (
    <div className="min-w-sm w-full max-w-3xl space-y-4">
      <h1 className="text-2xl font-semibold">Profile</h1>

      <Tabs defaultValue="info" className="w-full">
        <TabsList>
          <TabsTrigger value="info">Info</TabsTrigger>
          <TabsTrigger value="account">Account</TabsTrigger>

          <TabsTrigger value="settings">Settings</TabsTrigger>
        </TabsList>
        <TabsContent value="info">
          <Card className="p-4">
            <CardHeader>
              <div className="flex items-center gap-4">
                <Avatar size="lg">
                  <AvatarImage
                    src="/assets/avatar-placeholder.png"
                    alt="User"
                  />
                  <AvatarFallback>U</AvatarFallback>
                </Avatar>
                <div className="flex flex-col gap-1">
                  <div className="text-lg font-medium">
                    {user?.name || "User Name"}
                  </div>
                  <div className="text-sm text-muted-foreground">
                    {user?.email || "Email not available"}
                  </div>
                </div>
              </div>
              <CardAction>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button variant="outline">
                      <Edit2 className="h-4 w-4" />
                      Edit Profile
                    </Button>
                  </DialogTrigger>
                  <EditProfile />
                </Dialog>
              </CardAction>
            </CardHeader>

            <Separator />
            <CardContent>
              <h3 className="text-lg font-medium">QR Code</h3>
              <p className="text-sm text-muted-foreground">
                Scan this QR code to add me as a contact.
              </p>
              <QRCode uid={uid} />
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="account">
          <Card className="p-4">
            <CardHeader>
              <CardTitle>Account Details</CardTitle>
              <CardDescription>
                This section can include additional account details or settings.
              </CardDescription>
              <CardAction>
                <Button variant="destructive">Delete Account</Button>
              </CardAction>
            </CardHeader>
            <Separator />
            <CardContent>
              <form action="" className="flex flex-col gap-2">
                <p className="text-lg font-medium">Change Password</p>
                <Input placeholder="Email" />
                <Input
                  placeholder="Password old"
                  type="password"
                  className="mt-2"
                />
                <Input
                  placeholder="New Password"
                  type="password"
                  className="mt-2"
                />
                <Input
                  placeholder="Confirm Password"
                  type="password"
                  className="mt-2"
                />
                <Button variant="default" type="submit" className="mt-4">
                  Update Account
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="settings">
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
                  <Switch id="airplane-mode" />
                  <Label htmlFor="airplane-mode">Light/Dark</Label>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
