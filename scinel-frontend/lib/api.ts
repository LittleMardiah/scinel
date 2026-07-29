import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from 'axios';
import Cookies from 'js-cookie';
import { Toast } from '@/components/common/Toast';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const API_TIMEOUT = parseInt(process.env.NEXT_PUBLIC_API_TIMEOUT || '30000');

export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
  timeout: API_TIMEOUT,
});

// Request interceptor: Add JWT token to all requests...
apiClient.interceptors.request.use((config) => {
  const token = Cookies.get('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor: Handle errors & token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

    // Handle 401 - try refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = Cookies.get('refresh_token');
        if (!refreshToken) {
          Cookies.remove('access_token');
          Cookies.remove('refresh_token');
          window.location.href = '/login';
          return Promise.reject(error);
        }

        const response = await axios.post(
          `${API_BASE}/auth/refresh`,
          { refresh_token: refreshToken },
          { headers: { 'Content-Type': 'application/json' }, timeout: API_TIMEOUT }
        );

        const { access_token, refresh_token } = response.data.data;

        Cookies.set('access_token', access_token, {
          secure: true,
          sameSite: 'strict',
          expires: 1 / 24,
        });
        Cookies.set('refresh_token', refresh_token, {
          secure: true,
          sameSite: 'strict',
          expires: 7,
        });

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
        }
        return apiClient(originalRequest);
      } catch {
        Cookies.remove('access_token');
        Cookies.remove('refresh_token');
        window.location.href = '/login';
        return Promise.reject(error);
      }
    }

    // Handle 429 Too Many Requests
    if (error.response?.status === 429) {
      const retryAfter = error.response.headers['retry-after'] || '60';
      Toast.error(`Terlalu banyak permintaan. Coba lagi dalam ${retryAfter} detik.`);
      return Promise.reject({ ...error, retryAfter: parseInt(retryAfter) });
    }

    // Handle 403
    if (error.response?.status === 403) {
      Toast.error('Anda tidak berhak mengakses resource ini.');
      return Promise.reject(error);
    }

    // Handle 404
    if (error.response?.status === 404) {
      Toast.error('Resource tidak ditemukan.');
      return Promise.reject(error);
    }

    // Handle 500+
    if (error.response?.status && error.response.status >= 500) {
      const requestId = error.response.headers['x-request-id'] || 'unknown';
      Toast.error(`Terjadi kesalahan server (ID: ${requestId}). Tim kami sedang menangani.`);
      return Promise.reject({ ...error, requestId });
    }

    // Network errors
    if (!error.response) {
      Toast.error('Koneksi gagal. Periksa internet Anda dan coba lagi.');
      return Promise.reject(error);
    }

    Toast.error(error.message || 'Terjadi kesalahan. Coba lagi nanti.');
    return Promise.reject(error);
  }
);

export default apiClient;
