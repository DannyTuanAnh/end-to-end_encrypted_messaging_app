package main

import (
	"context"
	"log"

	p "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/package/cloud"
)

func main() {
	utils.LoadEnv()
	ctx := context.Background()

	// Giả lập sự kiện khi bạn vừa upload file "test.jpg" lên "raw-images"
	event := cloud.GCSEvent{
		Bucket: "chat-app-avt-images-raw",
		Name:   "c1164963-bfaa-4307-a8b6-7eadaaeb7149",
	}

	err := p.ProcessImage(ctx, event)
	if err != nil {
		log.Fatalf("Lỗi khi test: %v", err)
	}
	log.Println("Đã thông qua bước kiểm tra nội dung ảnh!")

	log.Println("Test hoàn tất thành công!")
}
