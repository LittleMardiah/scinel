'use client';

import { useState, FormEvent } from 'react';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';

interface RegisterFormProps {
  onSubmit: (email: string, password: string, fullName: string) => Promise<void>;
  isLoading?: boolean;
  error?: string | null;
}

export function RegisterForm({ onSubmit, isLoading, error }: RegisterFormProps) {
  const [formData, setFormData] = useState({ fullName: '', email: '', password: '', passwordConfirm: '' });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const errors: Record<string, string> = {};

    if (!formData.fullName) {
      errors.fullName = 'Nama lengkap diperlukan';
    }

    if (!formData.email) {
      errors.email = 'Email diperlukan';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = 'Format email tidak valid';
    }

    if (!formData.password) {
      errors.password = 'Kata sandi diperlukan';
    } else if (formData.password.length < 6) {
      errors.password = 'Kata sandi minimal 6 karakter';
    }

    if (formData.password !== formData.passwordConfirm) {
      errors.passwordConfirm = 'Kata sandi tidak cocok';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    await onSubmit(formData.email, formData.password, formData.fullName);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-200 rounded-md text-sm border border-red-200 dark:border-red-800">
          {error}
        </div>
      )}

      <div>
        <label htmlFor="fullName" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Nama Lengkap
        </label>
        <Input
          id="fullName"
          type="text"
          placeholder="Nama lengkap Anda"
          value={formData.fullName}
          onChange={(e) => setFormData({ ...formData, fullName: e.target.value })}
          error={formErrors.fullName}
        />
      </div>

      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Email
        </label>
        <Input
          id="email"
          type="email"
          placeholder="nama@perusahaan.com"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          error={formErrors.email}
        />
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Kata Sandi
        </label>
        <Input
          id="password"
          type="password"
          placeholder="Minimal 6 karakter"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          error={formErrors.password}
        />
      </div>

      <div>
        <label htmlFor="passwordConfirm" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Konfirmasi Kata Sandi
        </label>
        <Input
          id="passwordConfirm"
          type="password"
          placeholder="Masukkan ulang kata sandi"
          value={formData.passwordConfirm}
          onChange={(e) => setFormData({ ...formData, passwordConfirm: e.target.value })}
          error={formErrors.passwordConfirm}
        />
      </div>

      <Button type="submit" isLoading={isLoading} className="w-full">
        Daftar
      </Button>
    </form>
  );
}
