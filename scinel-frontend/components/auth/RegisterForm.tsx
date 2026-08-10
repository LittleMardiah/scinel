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
  const [formData, setFormData] = useState({ email: '', password: '', confirmPassword: '', fullName: '' });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});
  const [passwordStrength, setPasswordStrength] = useState<'weak' | 'medium' | 'strong' | null>(null);

  const validatePassword = (password: string) => {
    let strength: 'weak' | 'medium' | 'strong' = 'weak';
    if (password.length >= 8) strength = 'medium';
    if (password.length >= 12 && /[A-Z]/.test(password) && /[0-9]/.test(password) && /[^A-Za-z0-9]/.test(password)) {
      strength = 'strong';
    }
    setPasswordStrength(strength);
  };

  const validateForm = () => {
    const errors: Record<string, string> = {};

    if (!formData.fullName) errors.fullName = 'Nama lengkap diperlukan';
    if (!formData.email) {
      errors.email = 'Email diperlukan';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = 'Format email tidak valid';
    }
    if (!formData.password) {
      errors.password = 'Kata sandi diperlukan';
    } else if (formData.password.length < 8) {
      errors.password = 'Kata sandi minimal 8 karakter';
    }
    if (!formData.confirmPassword) {
      errors.confirmPassword = 'Konfirmasi kata sandi diperlukan';
    } else if (formData.confirmPassword !== formData.password) {
      errors.confirmPassword = 'Kata sandi tidak sama';
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
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nama Lengkap</label>
        <Input
          type="text"
          placeholder="Nama lengkap"
          value={formData.fullName}
          onChange={(e) => setFormData({ ...formData, fullName: e.target.value })}
          error={formErrors.fullName}
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email</label>
        <Input
          type="email"
          placeholder="nama@perusahaan.com"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          error={formErrors.email}
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Kata Sandi</label>
        <Input
          type="password"
          placeholder="Minimal 8 karakter"
          value={formData.password}
          onChange={(e) => {
            setFormData({ ...formData, password: e.target.value });
            validatePassword(e.target.value);
          }}
          error={formErrors.password}
        />
        {passwordStrength && (
          <div className="mt-1 flex items-center gap-2">
            <span className="text-xs text-gray-500">Kekuatan:</span>
            <span className={`text-xs font-medium ${
              passwordStrength === 'weak' ? 'text-red-500' :
              passwordStrength === 'medium' ? 'text-yellow-500' : 'text-green-500'
            }`}>
              {passwordStrength === 'weak' ? 'Lemah' :
               passwordStrength === 'medium' ? 'Sedang' : 'Kuat'}
            </span>
          </div>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Konfirmasi Kata Sandi</label>
        <Input
          type="password"
          placeholder="Ulangi kata sandi"
          value={formData.confirmPassword}
          onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
          error={formErrors.confirmPassword}
        />
      </div>

      <Button type="submit" isLoading={isLoading} className="w-full">
        Daftar
      </Button>
    </form>
  );
}
