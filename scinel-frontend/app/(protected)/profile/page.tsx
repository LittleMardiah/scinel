'use client';
import { useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Modal } from '@/components/common/Modal';
import { Toast } from '@/components/common/Toast';
import { apiClient } from '@/lib/api';
import { formatDistanceToNow, format } from 'date-fns';
import { id } from 'date-fns/locale';

export default function ProfilePage() {
  const { user, refreshUser } = useAuth();
  const [fullName, setFullName] = useState(user?.full_name || '');
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [pwModal, setPwModal] = useState(false);
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [pwError, setPwError] = useState('');

  const handleUpdateProfile = async () => {
    if (!fullName.trim()) return Toast.error('Nama tidak boleh kosong');
    setIsLoading(true);
    try {
      // TODO: Implement PUT /auth/me when backend ready. Mocking for now.
      console.log('[DEBUG] Update profile name to:', fullName);
      localStorage.setItem('mock_full_name', fullName);
      await refreshUser();
      Toast.success('✅ Profil berhasil diperbarui');
      setIsEditing(false);
    } catch {} finally { setIsLoading(false); }
  };

  const handleChangePassword = async () => {
    setPwError('');
    if (newPassword.length < 8) { setPwError('Minimal 8 karakter'); return; }
    if (newPassword !== confirmPassword) { setPwError('Password tidak cocok'); return; }
    try {
      await apiClient.put('/auth/change-password', { old_password: oldPassword, new_password: newPassword });
      Toast.success('✅ Kata sandi berhasil diubah');
      setPwModal(false); setOldPassword(''); setNewPassword(''); setConfirmPassword('');
    } catch (error: any) {
      if (error.response?.status === 401) setPwError('Kata sandi lama salah');
      else setPwError('Terjadi kesalahan, coba lagi.');
    }
  };

  if (!user) return <div className="p-6">Loading...</div>;

  return (
    <div className="p-6 max-w-3xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">👤 Profil Saya</h1>
      <Card>
        <CardHeader><CardTitle>Informasi Akun</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div><label className="text-sm text-gray-500">Email</label><p className="font-medium">{user.email}</p></div>
            <div><label className="text-sm text-gray-500">Role</label><p className="font-medium capitalize">{user.role}</p></div>
            <div><label className="text-sm text-gray-500">Terakhir Login</label><p className="font-medium">{user.last_login_at ? formatDistanceToNow(new Date(user.last_login_at), { addSuffix: true, locale: id }) : 'Belum pernah'}</p></div>
            <div><label className="text-sm text-gray-500">Akun Dibuat</label><p className="font-medium">{format(new Date(user.created_at), 'dd MMMM yyyy', { locale: id })}</p></div>
          </div>
          <div className="border-t pt-4">
            <label className="text-sm text-gray-500 block mb-1">Nama Lengkap</label>
            {isEditing ? (
              <div className="flex gap-2 items-center"><Input value={fullName} onChange={(e) => setFullName(e.target.value)} className="flex-1" /><Button onClick={handleUpdateProfile} isLoading={isLoading}>Simpan</Button><Button variant="ghost" onClick={() => { setIsEditing(false); setFullName(user.full_name || ''); }}>Batal</Button></div>
            ) : (
              <div className="flex items-center justify-between"><p className="font-medium">{user.full_name || '—'}</p><Button variant="ghost" onClick={() => setIsEditing(true)}>Edit</Button></div>
            )}
          </div>
          <Button variant="secondary" onClick={() => setPwModal(true)}>🔒 Ganti Kata Sandi</Button>
        </CardContent>
      </Card>
      <Modal isOpen={pwModal} onClose={() => setPwModal(false)} title="Ganti Kata Sandi">
        <div className="space-y-4">
          {pwError && <div className="p-2 bg-red-50 text-red-600 rounded text-sm">{pwError}</div>}
          <div><label className="block text-sm font-medium mb-1">Kata Sandi Lama</label><Input type="password" value={oldPassword} onChange={(e) => setOldPassword(e.target.value)} /></div>
          <div><label className="block text-sm font-medium mb-1">Kata Sandi Baru (min 8)</label><Input type="password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} /></div>
          <div><label className="block text-sm font-medium mb-1">Konfirmasi</label><Input type="password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} /></div>
          <div className="flex justify-end gap-2 pt-2"><Button variant="ghost" onClick={() => setPwModal(false)}>Batal</Button><Button onClick={handleChangePassword}>Ubah</Button></div>
        </div>
      </Modal>
    </div>
  );
}
