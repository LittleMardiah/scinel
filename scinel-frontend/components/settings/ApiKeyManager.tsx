'use client';
import { useState, useEffect } from 'react';
import { ApiKey, getApiKeys, createApiKey, revokeApiKey, rotateApiKey } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Modal } from '@/components/common/Modal';
import { Input } from '@/components/common/Input';
import { Select } from '@/components/common/Select';
import { Toast } from '@/components/common/Toast';
import { formatDistanceToNow } from 'date-fns';
import { id } from 'date-fns/locale';
import { RefreshCw, Trash2, Copy } from 'lucide-react';

export function ApiKeyManager() {
  const [keys, setKeys] = useState<ApiKey[]>([]); const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false); const [secretModal, setSecretModal] = useState<{ secret: string; prefix: string } | null>(null);
  const [formName, setFormName] = useState(''); const [formPermissions, setFormPermissions] = useState<string[]>(['inventory:read']); const [formExpiry, setFormExpiry] = useState('365');
  const permissionOptions = [
    { value: 'inventory:read', label: 'Baca Inventory' }, { value: 'inventory:write', label: 'Tulis Inventory' },
    { value: 'forecast:read', label: 'Baca Forecast' }, { value: 'forecast:write', label: 'Tulis Forecast' },
    { value: 'anomaly:read', label: 'Baca Anomaly' }, { value: 'webhook:manage', label: 'Kelola Webhook' },
  ];

  const loadKeys = async () => { setLoading(true); try { const data = await getApiKeys(); setKeys(data); } catch {} finally { setLoading(false); } };
  useEffect(() => { loadKeys(); }, []);

  const handleCreate = async () => {
    if (!formName.trim()) return Toast.error('Nama wajib diisi');
    if (formPermissions.length === 0) return Toast.error('Pilih minimal 1 permission');
    try {
      const result = await createApiKey({ name: formName, permissions: formPermissions, expires_in_days: parseInt(formExpiry) });
      Toast.success('✅ API Key berhasil dibuat!');
      setSecretModal({ secret: result.key_secret, prefix: result.key_prefix });
      setModalOpen(false); setFormName(''); setFormPermissions(['inventory:read']); setFormExpiry('365');
      await loadKeys();
    } catch {}
  };

  const handleRevoke = async (id: string, name: string) => {
    if (!confirm(`Yakin mencabut "${name}"?`)) return;
    try { await revokeApiKey(id); Toast.success(`✅ "${name}" dicabut`); await loadKeys(); } catch {}
  };
  const handleRotate = async (id: string, name: string) => {
    if (!confirm(`Yakin mengganti secret "${name}"?`)) return;
    try { const result = await rotateApiKey(id); Toast.success('🔄 Secret diganti!'); setSecretModal({ secret: result.key_secret, prefix: 'sk_' }); await loadKeys(); } catch {}
  };
  const copyToClipboard = (text: string) => { navigator.clipboard.writeText(text); Toast.info('📋 Tersalin!'); };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-lg">🔑 API Keys</CardTitle>
        <Button onClick={() => setModalOpen(true)}>Buat Baru</Button>
      </CardHeader>
      <CardContent>
        {loading ? <p>Loading...</p> : keys.length === 0 ? <p className="text-gray-500">Belum ada API Key.</p> : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="text-gray-600 dark:text-gray-400 border-b">
                <tr><th className="text-left py-2">Prefix</th><th className="text-left py-2">Nama</th><th className="text-left py-2">Izin</th><th className="text-left py-2">Status</th><th className="text-left py-2">Terakhir</th><th className="text-left py-2">Aksi</th></tr>
              </thead>
              <tbody>
                {keys.map((key) => (
                  <tr key={key.id} className="border-t border-gray-100 dark:border-gray-800">
                    <td className="py-2 font-mono text-xs">{key.key_prefix}...</td>
                    <td className="py-2">{key.name}</td>
                    <td className="py-2 text-xs">{key.permissions.join(', ')}</td>
                    <td className="py-2"><span className={`px-2 py-0.5 rounded-full text-xs ${key.is_active ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300' : 'bg-gray-100 text-gray-500'}`}>{key.is_active ? 'Aktif' : 'Nonaktif'}</span></td>
                    <td className="py-2 text-xs text-gray-500">{key.last_used_at ? formatDistanceToNow(new Date(key.last_used_at), { addSuffix: true, locale: id }) : 'Belum dipakai'}</td>
                    <td className="py-2 flex gap-1">
                      <Button size="sm" variant="ghost" onClick={() => handleRotate(key.id, key.name)}><RefreshCw className="w-4 h-4" /></Button>
                      <Button size="sm" variant="ghost" onClick={() => handleRevoke(key.id, key.name)}><Trash2 className="w-4 h-4 text-red-500" /></Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        <Modal isOpen={modalOpen} onClose={() => setModalOpen(false)} title="Buat API Key Baru">
          <div className="space-y-4">
            <div><label className="block text-sm font-medium mb-1">Nama</label><Input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder="Contoh: ERP Integration" /></div>
            <div><label className="block text-sm font-medium mb-1">Izin</label><Select multiple value={formPermissions} onChange={(val) => setFormPermissions(Array.isArray(val) ? val : [val])} options={permissionOptions} /></div>
            <div><label className="block text-sm font-medium mb-1">Kedaluwarsa</label><Select value={formExpiry} onChange={(val) => setFormExpiry(val as string)} options={[{ value: '30', label: '30 hari' }, { value: '90', label: '90 hari' }, { value: '365', label: '1 tahun' }, { value: '0', label: 'Tidak pernah' }]} /></div>
            <div className="flex justify-end gap-2 pt-2"><Button variant="ghost" onClick={() => setModalOpen(false)}>Batal</Button><Button onClick={handleCreate}>Buat</Button></div>
          </div>
        </Modal>
        <Modal isOpen={!!secretModal} onClose={() => setSecretModal(null)} title="⚠️ Simpan Secret Ini!">
          <div className="space-y-4">
            <p className="text-red-600 dark:text-red-400 font-medium">Secret hanya muncul SEKALI!</p>
            <div className="bg-gray-100 dark:bg-gray-800 p-3 rounded font-mono text-sm break-all flex items-center justify-between"><span>{secretModal?.secret}</span><Button size="sm" variant="ghost" onClick={() => copyToClipboard(secretModal?.secret || '')}><Copy className="w-4 h-4" /></Button></div>
            <Button onClick={() => setSecretModal(null)} className="w-full">Saya sudah salin</Button>
          </div>
        </Modal>
      </CardContent>
    </Card>
  );
}
