'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { RegisterForm } from '@/components/auth/RegisterForm';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Toast } from '@/components/common/Toast';

export default function RegisterPage() {
  const router = useRouter();
  const { register } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleRegister = async (email: string, password: string, fullName: string) => {
    setIsLoading(true);
    setError(null);
    try {
      await register(email, password, fullName);
      Toast.success('Akun berhasil dibuat! Selamat datang di ScIntel.');
      router.push('/dashboard');
    } catch (err: any) {
      const message = err.response?.data?.message || err.message || 'Registrasi gagal. Coba lagi.';
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
            <CardTitle className="text-2xl mb-2">Daftar ScIntel</CardTitle>
            <p className="text-gray-600 dark:text-gray-400">
              Buat akun untuk mulai
            </p>
          </CardHeader>

          <CardContent>
            <RegisterForm
              onSubmit={handleRegister}
              isLoading={isLoading}
              error={error}
            />

            <div className="mt-6 text-center">
              <p className="text-gray-600 dark:text-gray-400 text-sm">
                Sudah punya akun?{' '}
                <Link
                  href="/login"
                  className="text-blue-600 hover:text-blue-700 font-medium"
                >
                  Masuk
                </Link>
              </p>
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
