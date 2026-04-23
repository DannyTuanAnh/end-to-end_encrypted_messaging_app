import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuthContext } from "@/context/AuthContext";
import { Spinner } from "@/components/ui/spinner";
import Info from "@/pages/profile/Info";
import Account from "@/pages/profile/Account";
import Setting from "@/pages/profile/Setting";

export default function Profile() {
  const { user } = useAuthContext();
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
          {user === null ? (
            <div className="flex items-center justify-center h-48">
              <Spinner />
            </div>
          ) : (
            <Info
              uid={user.uid}
              name={user.name || undefined}
              email={user.email || undefined}
            />
          )}
        </TabsContent>
        <TabsContent value="account">
          <Account
            handleDeleteAccount={() => {
              console.log("Delete account");
            }}
            handleUpdateAccount={(data) => {
              console.log("Update account", data);
            }}
          />
        </TabsContent>

        <TabsContent value="settings">
          <Setting />
        </TabsContent>
      </Tabs>
    </div>
  );
}
