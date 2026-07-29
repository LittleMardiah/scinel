'use client';

import React, { useState } from 'react';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';

interface ResetPasswordFormProps {
  onSubmit: (newPassword: string) => Promise<void>;
  isLoading?: boolean;
  error?: string | null;
}

export const ResetPasswordForm: React.FC<ResetPasswordFormProps> = ({
  onSubmit,
  isLoading,
  error,
}) => {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});
  const [passwordStrength, setPasswordStrength] = useState<'weak' | 'medium' | 'strong' | ''>('');

  const checkPasswordStrength = (pass: string): 'weak' | 'medium' | 'strong' | '' => {
    if (!pass) return '';
    let score = 0;
    if (pass.length >= 8) score++;
    if (/[a-z]/.test(pass)) score++;
    if (/[A-Z]/.test(pass)) score++;
    if (/[0-9]/.test(pass)) score++;
    if (/[^a-zA-Z0-9]/.test(pass)) score++;
    if (score <= 2) return 'weak';
    if (score <= 4) return 'medium';
    return 'strong';
  };

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};

    if (!password) {
      errors.password = 'Kata sandi baru diperlukan';
    } else if (password.length < 8) {
      errors.password = 'Minimal 8 karakter';
    }

    if (!confirmPassword) {
      errors.confirmPassword = 'Konfirmasi kata sandi diperlukan';
    } else if (password && confirmPassword && password !== confirmPassword) {
      errors.confirmPassword = 'Kata sandi tidak cocok';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
    setPasswordStrength(checkPasswordStrength(value));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;
    onSubmit(password);
  };

  const strengthColors = {
    weak: 'bg-red-500',
    medium: 'bg-yellow-500',
    strong: 'bg-green-500',
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-200 rounded-md text-sm border border-red-200 dark:border-red-800">
          {error}
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Kata Sandi Baru
        </label>
        <Input
          type="password"
          placeholder="Minimal 8 karakter, huruf & angka"
          value={password}
          onChange={handlePasswordChange}
          error={formErrors.password}
          disabled={isLoading}
          autoComplete="new-password"
        />
        {password && passwordStrength && (
          <div className="mt-1">
            <div className="flex items-center gap-2">
              <div className="h-1.5 flex-1 rounded-full bg-gray-200 dark:bg-gray-700 overflow-hidden">
                <div
                  className={`h-full rounded-full transition-all ${
                    passwordStrength === 'weak' ? 'w-1/3' : passwordStrength === 'medium' ? 'w-2/3' : 'w-full'
                  } ${strengthColors[passwordStrength]}`}
                />
              </div>
              <span className="text-xs capitalize text-gray-500 dark:text-gray-400">{passwordStrength}</span>
            </div>
          </div>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Konfirmasi Kata Sandi Baru
        </label>
        <Input
          type="password"
          placeholder="Ulangi kata sandi baru"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          error={formErrors.confirmPassword}
          disabled={isLoading}
          autoComplete="new-password"
        />
      </div>

      <Button
        type="submit"
        className="w-full"
        isLoading={isLoading}
        disabled={isLoading || Object.keys(formErrors).length > 0}
      >
        {isLoading ? 'Memproses...' : 'Reset Kata Sandi'}
      </Button>
    </form>
  );
};
