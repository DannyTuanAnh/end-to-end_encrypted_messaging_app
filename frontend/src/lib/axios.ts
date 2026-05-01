import axios, { AxiosError, AxiosResponse } from 'axios';

// Base URL từ tài liệu (Hoặc lấy từ biến môi trường Vite nếu bạn có cấu hình)
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://chat-app-ta.duckdns.org/api/v1';

// Khởi tạo instance Axios với các cài đặt mặc định
export const api = axios.create({
  baseURL: API_BASE_URL,
  // RẤT QUAN TRỌNG: Cần thiết để gửi và nhận Cookie (session_id, device_id)
  withCredentials: true, 
  headers: {
    'Content-Type': 'application/json',
    // Nếu backend bắt buộc dùng API Key ở tất cả các route, hãy bỏ comment dòng dưới
    // 'X-Api-Key': import.meta.env.VITE_API_KEY || 'your-api-key',
  },
  timeout: 15000, // Timeout sau 15 giây
});

// Interceptor cho Request: Xử lý trước khi gửi request đi
api.interceptors.request.use(
  (config) => {
    // Có thể thực hiện log, chặn request, hoặc thêm logic khác tại đây
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor cho Response: Xử lý response trả về trước khi tới component
api.interceptors.response.use(
  (response: AxiosResponse) => {
    // Backend trả về format chuẩn { success: true, data: ..., message: ... }
    // Ta có thể bóc tách luôn `response.data` ở đây để khi gọi API bên ngoài code sẽ gọn hơn.
    return response.data;
  },
  (error: AxiosError) => {
    // Xử lý lỗi toàn cục từ Backend
    if (error.response) {
      const status = error.response.status;
      const data: any = error.response.data;

      // Log lỗi ra console để dễ debug
      console.error(`[API Error] Status: ${status}`, data);

      if (status === 401) {
        // Mã 401 Unauthorized: Session hết hạn hoặc không hợp lệ.
        // Gợi ý: Có thể dispatch event báo hết hạn phiên đăng nhập hoặc redirect user về trang đăng nhập.
        // window.location.href = '/login'; 
      }
      
      if (status === 403) {
        // Mã 403 Forbidden: Thường là do API Key bị sai hoặc không có quyền.
        console.error('Lỗi phân quyền hoặc thiếu API Key!');
      }

    } else if (error.request) {
      console.error('[API Error] Không thể kết nối tới server:', error.request);
    } else {
      console.error('[API Error] Lỗi cấu hình:', error.message);
    }

    return Promise.reject(error);
  }
);

export default api;
