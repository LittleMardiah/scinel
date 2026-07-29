'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { ForgotPasswordForm } from '@/components/auth/ForgotPasswordForm';
import { Toast } from '@/components/common/Toast';
import apiClient from '@/lib/api';

export default function ForgotPasswordPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleSubmit = async (email: string) => {
    setIsLoading(true);
    setError(null);
    setSuccess(null);
    try {
      await apiClient.post('/auth/password-reset', { email });
      setSuccess('Link reset kata sandi telah dikirim ke email Anda. Cek inbox/spam.');
      Toast.success('Email reset terkirim!');
    } catch (err: any) {
      const message = err.response?.data?.message || 'Gagal mengirim link reset. Coba lagi.';
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
            <CardTitle className="text-2xl mb-2">Lupa Kata Sandi</CardTitle>
            <p className="text-gray-600 dark:text-gray-400">
              Masukkan email Anda, kami akan kirim link reset.
            </p>
          </CardHeader>

          <CardContent>
            <ForgotPasswordForm
              onSubmit={handleSubmit}
              isLoading={isLoading}
              error={error}
              success={success}
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
