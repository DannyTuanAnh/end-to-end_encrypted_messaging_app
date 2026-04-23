import React from "react";
import { QRCodeCanvas } from "qrcode.react";

type Props = {
  uid: string;
};

const QRCode: React.FC<Props> = ({ uid }) => {
  return (
    <div>
      <QRCodeCanvas
        value={uid} // dữ liệu QR (uid)
        size={200} // kích thước
        level="H" // độ chịu lỗi
        includeMargin // padding
      />
    </div>
  );
};

export default QRCode;
