'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { LoginForm } from '@/components/auth/LoginForm';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Toast } from '@/components/common/Toast';

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);
    try {
      await login(email, password);
      Toast.success('Berhasil masuk! Membawa Anda ke dashboard...');
      router.push('/dashboard');
    } catch (err: any) {
      const message =
        err.response?.data?.message ||
        err.message ||
        'Login gagal. Coba lagi.';
      setError(message);
      Toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
      <div className="w-full max-w-md">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-2xl mb-2">Masuk ke ScIntel</CardTitle>
            <p className="text-gray-600 dark:text-gray-400">
              Kelola supply chain Anda dengan cerdas
            </p>
          </CardHeader>

          <CardContent>
            <LoginForm
              onSubmit={handleLogin}
              isLoading={isLoading}
              error={error}
            />

            <div className="mt-6 space-y-3">
              <div className="text-center">
                <Link
                  href="/forgot-password"
                  className="text-blue-600 hover:text-blue-700 text-sm font-medium"
                >
                  Lupa kata sandi?
                </Link>
              </div>
              <div className="text-center">
                <p className="text-gray-600 dark:text-gray-400 text-sm">
                  Belum punya akun?{' '}
                  <Link
                    href="/register"
                    className="text-blue-600 hover:text-blue-700 font-medium"
                  >
                    Daftar sekarang
                  </Link>
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="text-center mt-6 text-xs text-gray-500">
          <p>ScIntel v1.0.0 • © 2026 Supply Chain Intelligence</p>
        </div>
      </div>
    </div>
  );
}
