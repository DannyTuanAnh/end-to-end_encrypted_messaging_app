import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { updatePasswordSchema } from "@/lib/schema";
import type { UpdatePasswordForm } from "@/lib/schema";
import { Trash2 } from "lucide-react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
type AccountProps = {
  handleDeleteAccount: () => void;
  handleUpdateAccount: (data: {
    email: string;
    oldPassword: string;
    newPassword: string;
    confirmPassword: string;
  }) => void;
};
export default function Account({
  handleDeleteAccount,
  handleUpdateAccount,
}: AccountProps) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<UpdatePasswordForm>({
    resolver: zodResolver(updatePasswordSchema),
    mode: "onChange",
  });
  const onSubmit = (data: UpdatePasswordForm) => {
    const formattedData = {
      email: data.email,
      oldPassword: data.passwordOld,
      newPassword: data.passwordNew,
      confirmPassword: data.passwordConfirm,
    };
    handleUpdateAccount(formattedData);
  };
  return (
    <Card className="p-4">
      <CardHeader>
        <CardTitle>Account Details</CardTitle>
        <CardDescription>
          This section can include additional account details or settings.
        </CardDescription>
        <CardAction>
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="destructive">
                <Trash2 className="h-4 w-4" />
                Delete Account
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                <AlertDialogDescription>
                  This action cannot be undone. This will permanently delete
                  your account and remove your data from our servers.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  variant="destructive"
                  onClick={handleDeleteAccount}
                >
                  Delete Account
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </CardAction>
      </CardHeader>
      <Separator />
      <CardContent>
        <form
          action=""
          className="flex flex-col gap-2"
          onSubmit={handleSubmit(onSubmit)}
        >
          <p className="text-lg font-medium">Change Password</p>
          <Input placeholder="Email" type="email" {...register("email")} />
          {errors.email && (
            <p className="text-red-500">{errors.email.message}</p>
          )}
          <Input
            placeholder="Password old"
            type="password"
            className="mt-2"
            {...register("passwordOld")}
          />
          {errors.passwordOld && (
            <p className="text-red-500">{errors.passwordOld.message}</p>
          )}
          <Input
            placeholder="New Password"
            type="password"
            className="mt-2"
            {...register("passwordNew")}
          />
          {errors.passwordNew && (
            <p className="text-red-500">{errors.passwordNew.message}</p>
          )}
          <Input
            placeholder="Confirm Password"
            type="password"
            className="mt-2"
            {...register("passwordConfirm")}
          />
          {errors.passwordConfirm && (
            <p className="text-red-500">{errors.passwordConfirm.message}</p>
          )}
          <Button
            variant="default"
            type="submit"
            className="mt-4"
            disabled={isSubmitting}
          >
            Update Account
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
