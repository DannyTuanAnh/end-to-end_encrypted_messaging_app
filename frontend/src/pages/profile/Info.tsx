import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Dialog, DialogTrigger } from "@/components/ui/dialog";
import { EditProfile } from "@/components/profile/EditProfile";
import QRCode from "@/components/common/QRCode";
import { Button } from "@/components/ui/button";
import { Edit2 } from "lucide-react";
import { Separator } from "@/components/ui/separator";
type Props = {
  uid: string;
  name?: string;
  email?: string;
};
export default function Info({ uid, name, email }: Props) {
  return (
    <Card className="p-4">
      <CardHeader>
        <div className="flex items-center gap-4">
          <Avatar size="lg">
            <AvatarImage src="/assets/avatar-placeholder.png" alt="User" />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
          <div className="flex flex-col gap-1">
            <div className="text-lg font-medium">{name || "User Name"}</div>
            <div className="text-sm text-muted-foreground">
              {email || "Email not available"}
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
  );
}
