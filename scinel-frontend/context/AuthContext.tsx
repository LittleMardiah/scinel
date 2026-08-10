'use client';

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  ReactNode,
} from 'react';
import Cookies from 'js-cookie';
import { apiClient } from '@/lib/api';

export interface User {
  id: string;
  email: string;
  full_name?: string;
  role: 'admin' | 'manager' | 'analyst' | 'viewer';
  status: 'active' | 'suspended' | 'pending_email_verification';
  created_at: string;
  last_login_at?: string;
}

export interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const decodeToken = useCallback((token: string): User | null => {
    try {
      const parts = token.split('.');
      if (parts.length !== 3) return null;
      const payload = JSON.parse(atob(parts[1]));
      return {
        id: payload.user_id,
        email: payload.email,
        role: payload.role,
        status: payload.status || 'active',
        created_at: new Date(payload.iat * 1000).toISOString(),
      };
    } catch {
      return null;
    }
  }, []);

  useEffect(() => {
    const token = Cookies.get('access_token');
    if (token) {
      const decodedUser = decodeToken(token);
      if (decodedUser) {
        setUser(decodedUser);
      } else {
        Cookies.remove('access_token');
        Cookies.remove('refresh_token');
      }
    }
    setIsLoading(false);
  }, [decodeToken]);

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const response = await apiClient.post('/auth/login', { email, password });
      const { access_token, refresh_token, user: userData } = response.data.data;
      Cookies.set('access_token', access_token, { secure: true, sameSite: 'strict', expires: 1 / 24 });
      Cookies.set('refresh_token', refresh_token, { secure: true, sameSite: 'strict', expires: 7 });
      setUser(userData);
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (email: string, password: string, fullName: string) => {
    setIsLoading(true);
    try {
      const response = await apiClient.post('/auth/register', { email, password, full_name: fullName });
      const { access_token, refresh_token, user: userData } = response.data.data;
      Cookies.set('access_token', access_token, { secure: true, sameSite: 'strict', expires: 1 / 24 });
      Cookies.set('refresh_token', refresh_token, { secure: true, sameSite: 'strict', expires: 7 });
      setUser(userData);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      const refreshToken = Cookies.get('refresh_token');
      if (refreshToken) {
        await apiClient.post('/auth/logout', { refresh_token: refreshToken });
      }
    } catch (error) {
      console.warn('Logout error:', error);
    } finally {
      Cookies.remove('access_token');
      Cookies.remove('refresh_token');
      setUser(null);
    }
  };

  const refreshUser = async () => {
    const token = Cookies.get('access_token');
    if (token) {
      const decodedUser = decodeToken(token);
      setUser(decodedUser);
    }
  };

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, login, register, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
}
