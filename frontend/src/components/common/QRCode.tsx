import React from "react";
import { QRCodeCanvas } from "qrcode.react";

type Props = {
  uid: string | undefined;
};

const QRCode: React.FC<Props> = ({ uid }) => {
  return (
    <div>
      {uid ? (
        <QRCodeCanvas
          value={uid} // dữ liệu QR (uid)
          size={200} // kích thước
          level="H" // độ chịu lỗi
          includeMargin // padding
        />
      ) : (
        <p className="text-sm text-muted-foreground">
          No user ID available to generate QR code.
        </p>
      )}
    </div>
  );
};

export default QRCode;
