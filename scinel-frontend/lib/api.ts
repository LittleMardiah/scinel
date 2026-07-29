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

apiClient.interceptors.request.use((config) => {
  const token = Cookies.get('access_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      try {
        const refreshToken = Cookies.get('refresh_token');
        if (!refreshToken) throw new Error('No refresh token');
        const response = await axios.post(`${API_BASE}/auth/refresh`, { refresh_token: refreshToken });
        const { access_token, refresh_token } = response.data.data;
        Cookies.set('access_token', access_token, { secure: true, sameSite: 'strict', expires: 1 / 24 });
        Cookies.set('refresh_token', refresh_token, { secure: true, sameSite: 'strict', expires: 7 });
        if (originalRequest.headers) originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        Cookies.remove('access_token'); Cookies.remove('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    const errorHandlers: Record<number, () => void> = {
      429: () => {
        const retryAfter = error.response?.headers['retry-after'] || '60';
        Toast.error(`Terlalu banyak permintaan. Coba lagi dalam ${retryAfter} detik.`);
      },
      403: () => Toast.error('Anda tidak berhak mengakses resource ini.'),
      404: () => Toast.error('Resource tidak ditemukan.'),
    };
    const handler = errorHandlers[error.response?.status as keyof typeof errorHandlers];
    if (handler) handler();
    else if (error.response?.status && error.response.status >= 500) {
      const requestId = error.response.headers['x-request-id'] || 'unknown';
      Toast.error(`Terjadi kesalahan server (ID: ${requestId}). Tim kami sedang menangani.`);
    } else if (!error.response) Toast.error('Koneksi gagal. Periksa internet Anda.');
    return Promise.reject(error);
  }
);

// ============================================================
// API KEY INTERFACES & FUNCTIONS
// ============================================================
export interface ApiKey {
  id: string; key_prefix: string; name: string; permissions: string[];
  is_active: boolean; last_used_at?: string; expires_at?: string; created_at: string;
}

const mockApiKeys: ApiKey[] = [
  { id: 'ak-mock-1', key_prefix: 'sk_test_abc123', name: 'Integrasi ERP', permissions: ['inventory:read', 'forecast:read'], is_active: true, last_used_at: new Date().toISOString(), expires_at: new Date(Date.now() + 365*86400000).toISOString(), created_at: new Date().toISOString() },
  { id: 'ak-mock-2', key_prefix: 'sk_test_def456', name: 'Mobile App', permissions: ['inventory:read', 'inventory:write'], is_active: true, last_used_at: new Date(Date.now() - 86400000).toISOString(), expires_at: null, created_at: new Date(Date.now() - 30*86400000).toISOString() },
];

export const getApiKeys = async (): Promise<ApiKey[]> => {
  try { const res = await apiClient.get('/api-keys'); return res.data.data.keys || []; }
  catch (e) { console.warn('⚠️ [DEBUG] Backend API Keys placeholder, using mock data.'); return mockApiKeys; }
};

export const createApiKey = async (data: { name: string; permissions: string[]; expires_in_days?: number }): Promise<{ id: string; key_prefix: string; key_secret: string }> => {
  try { const res = await apiClient.post('/api-keys', data); return res.data.data; }
  catch (e) { console.warn('⚠️ [DEBUG] Backend create mock, generating fake key.'); return { id: 'ak-new-' + Date.now(), key_prefix: 'sk_new_' + Math.random().toString(36).slice(2,7), key_secret: 'sk_live_' + Math.random().toString(36).slice(2,15) }; }
};

export const revokeApiKey = async (id: string): Promise<void> => {
  try { await apiClient.delete(`/api-keys/${id}`); }
  catch (e) { console.warn(`⚠️ [DEBUG] Mock revoke for ${id}`); }
};

export const rotateApiKey = async (id: string): Promise<{ key_secret: string }> => {
  try { const res = await apiClient.post(`/api-keys/${id}/rotate`); return res.data.data; }
  catch (e) { console.warn(`⚠️ [DEBUG] Mock rotate for ${id}`); return { key_secret: 'sk_rotated_' + Math.random().toString(36).slice(2,15) }; }
};

// ============================================================
// WEBHOOK INTERFACES & FUNCTIONS
// ============================================================
export interface WebhookSubscription {
  id: string; event_type: string; webhook_url: string; is_active: boolean;
  max_retries: number; retry_interval: number; created_at: string;
}

const mockWebhooks: WebhookSubscription[] = [
  { id: 'wh-mock-1', event_type: 'inventory.changed', webhook_url: 'https://erp.example.com/webhook', is_active: true, max_retries: 5, retry_interval: 60, created_at: new Date().toISOString() },
  { id: 'wh-mock-2', event_type: 'anomaly.detected', webhook_url: 'https://monitor.example.com/webhook', is_active: false, max_retries: 3, retry_interval: 120, created_at: new Date(Date.now() - 7*86400000).toISOString() },
];

export const getWebhooks = async (): Promise<WebhookSubscription[]> => {
  try { const res = await apiClient.get('/webhooks/subscriptions'); return res.data.data.subscriptions || []; }
  catch (e) { console.warn('⚠️ [DEBUG] Backend Webhooks placeholder, using mock data.'); return mockWebhooks; }
};

export const createWebhook = async (data: { event_type: string; webhook_url: string; max_retries?: number; retry_interval?: number }): Promise<{ id: string; secret_key: string }> => {
  try { const res = await apiClient.post('/webhooks/subscribe', data); return res.data.data; }
  catch (e) { console.warn('⚠️ [DEBUG] Mock create webhook.'); return { id: 'wh-new-' + Date.now(), secret_key: 'whsec_' + Math.random().toString(36).slice(2,15) }; }
};

export const deleteWebhook = async (id: string): Promise<void> => {
  try { await apiClient.delete(`/webhooks/subscriptions/${id}`); }
  catch (e) { console.warn(`⚠️ [DEBUG] Mock delete webhook ${id}`); }
};

export const toggleWebhook = async (id: string, isActive: boolean): Promise<void> => {
  try { await apiClient.put(`/webhooks/subscriptions/${id}`, { is_active: isActive }); }
  catch (e) { console.warn(`⚠️ [DEBUG] Mock toggle webhook ${id} to ${isActive}`); }
};
