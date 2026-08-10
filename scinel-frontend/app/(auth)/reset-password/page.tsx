'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { ResetPasswordForm } from '@/components/auth/ResetPasswordForm';
import { Toast } from '@/components/common/Toast';
import { apiClient } from '@/lib/api';

export default function ResetPasswordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get('token');

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isTokenValid, setIsTokenValid] = useState<boolean | null>(null);

  useEffect(() => {
    if (!token) {
      setIsTokenValid(false);
      setError('Token reset tidak ditemukan. Link mungkin tidak valid.');
    } else {
      setIsTokenValid(true);
    }
  }, [token]);

  const handleSubmit = async (newPassword: string) => {
    if (!token) {
      setError('Token tidak valid.');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      await apiClient.put(`/auth/password-reset/${token}`, { new_password: newPassword });
      Toast.success('Kata sandi berhasil direset! Silakan login.');
      router.push('/login');
    } catch (err: any) {
      const message = err.response?.data?.message || 'Gagal mereset kata sandi. Coba lagi.';
      setError(message);
      Toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  if (isTokenValid === false) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
        <div className="w-full max-w-md">
          <Card>
            <CardHeader className="text-center">
              <CardTitle className="text-2xl mb-2">Link Tidak Valid</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-red-600 dark:text-red-400 text-center">
                {error || 'Token reset tidak ditemukan atau sudah kadaluwarsa.'}
              </p>
              <div className="mt-6 text-center">
                <Link
                  href="/forgot-password"
                  className="text-blue-600 hover:text-blue-700 text-sm font-medium"
                >
                  Minta link baru
                </Link>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
      <div className="w-full max-w-md">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-2xl mb-2">Reset Kata Sandi</CardTitle>
            <p className="text-gray-600 dark:text-gray-400">
              Masukkan kata sandi baru Anda.
            </p>
          </CardHeader>

          <CardContent>
            <ResetPasswordForm
              onSubmit={handleSubmit}
              isLoading={isLoading}
              error={error}
            />

            <div className="mt-6 text-center">
              <Link
                href="/login"
                className="text-blue-600 hover:text-blue-700 text-sm font-medium"
              >
                Kembali ke login
              </Link>
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
