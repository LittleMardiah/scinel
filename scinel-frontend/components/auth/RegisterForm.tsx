'use client';

import React, { useState } from 'react';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';

interface RegisterFormProps {
  onSubmit: (email: string, password: string, fullName: string) => Promise<void>;
  isLoading?: boolean;
  error?: string | null;
}

export const RegisterForm: React.FC<RegisterFormProps> = ({ onSubmit, isLoading, error }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [fullName, setFullName] = useState('');
  const [passwordError, setPasswordError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      setPasswordError('Passwords do not match');
      return;
    }
    setPasswordError('');
    onSubmit(email, password, fullName);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && <div className="rounded-md bg-red-50 p-3 text-sm text-red-700 dark:bg-red-900/20 dark:text-red-400">{error}</div>}
      <Input type="text" placeholder="Full Name" value={fullName} onChange={(e) => setFullName(e.target.value)} required disabled={isLoading} />
      <Input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required disabled={isLoading} />
      <Input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} required disabled={isLoading} />
      <Input type="password" placeholder="Confirm Password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required disabled={isLoading} error={passwordError} />
      <Button type="submit" className="w-full" isLoading={isLoading} disabled={isLoading}>
        {isLoading ? 'Loading...' : 'Register'}
      </Button>
    </form>
  );
};
