'use client';
import { useState, useEffect } from 'react';
import { WebhookSubscription, getWebhooks, createWebhook, deleteWebhook, toggleWebhook } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Modal } from '@/components/common/Modal';
import { Input } from '@/components/common/Input';
import { Select } from '@/components/common/Select';
import { Toast } from '@/components/common/Toast';
import { formatDistanceToNow } from 'date-fns';
import { id } from 'date-fns/locale';
import { Trash2, Eye } from 'lucide-react';

// Simple Switch component inline
const Switch = ({ checked, onChange }: { checked: boolean; onChange: () => void }) => (
  <button onClick={onChange} className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${checked ? 'bg-blue-600' : 'bg-gray-300 dark:bg-gray-600'}`}>
    <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${checked ? 'translate-x-6' : 'translate-x-1'}`} />
  </button>
);

export function WebhookManager() {
  const [webhooks, setWebhooks] = useState<WebhookSubscription[]>([]); const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [formEvent, setFormEvent] = useState('inventory.changed'); const [formUrl, setFormUrl] = useState(''); const [formRetries, setFormRetries] = useState('5'); const [formInterval, setFormInterval] = useState('60');
  const eventOptions = [{ value: 'inventory.changed', label: '📦 Inventory Berubah' }, { value: 'forecast.generated', label: '📊 Forecast Baru' }, { value: 'anomaly.detected', label: '🚨 Anomali' }];

  const loadWebhooks = async () => { setLoading(true); try { const data = await getWebhooks(); setWebhooks(data); } catch {} finally { setLoading(false); } };
  useEffect(() => { loadWebhooks(); }, []);

  const handleCreate = async () => {
    if (!formUrl.trim() || !formUrl.startsWith('https://')) return Toast.error('URL harus https://');
    try {
      const result = await createWebhook({ event_type: formEvent, webhook_url: formUrl, max_retries: parseInt(formRetries), retry_interval: parseInt(formInterval) });
      Toast.success('✅ Webhook terdaftar!'); Toast.info(`🔑 Secret: ${result.secret_key}`);
      setModalOpen(false); setFormUrl(''); await loadWebhooks();
    } catch {}
  };
  const handleDelete = async (id: string) => { if (!confirm('Hapus webhook?')) return; try { await deleteWebhook(id); Toast.success('Terhapus'); await loadWebhooks(); } catch {} };
  const handleToggle = async (id: string, current: boolean) => { try { await toggleWebhook(id, !current); Toast.success(`Webhook ${!current ? 'diaktifkan' : 'dinonaktifkan'}`); await loadWebhooks(); } catch {} };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-lg">🔗 Webhook Subscriptions</CardTitle>
        <Button onClick={() => setModalOpen(true)}>Subscribe</Button>
      </CardHeader>
      <CardContent>
        {loading ? <p>Loading...</p> : webhooks.length === 0 ? <p className="text-gray-500">Belum ada webhook.</p> : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="text-gray-600 dark:text-gray-400 border-b">
                <tr><th className="text-left py-2">Event</th><th className="text-left py-2">URL</th><th className="text-left py-2">Status</th><th className="text-left py-2">Dibuat</th><th className="text-left py-2">Aksi</th></tr>
              </thead>
              <tbody>
                {webhooks.map((wh) => (
                  <tr key={wh.id} className="border-t border-gray-100 dark:border-gray-800">
                    <td className="py-2">{wh.event_type}</td>
                    <td className="py-2 text-xs truncate max-w-[150px]">{wh.webhook_url}</td>
                    <td className="py-2"><Switch checked={wh.is_active} onChange={() => handleToggle(wh.id, wh.is_active)} /></td>
                    <td className="py-2 text-xs text-gray-500">{formatDistanceToNow(new Date(wh.created_at), { addSuffix: true, locale: id })}</td>
                    <td className="py-2 flex gap-1">
                      <Button size="sm" variant="ghost" onClick={() => window.open(`/webhooks/deliveries?sub=${wh.id}`, '_blank')}><Eye className="w-4 h-4" /></Button>
                      <Button size="sm" variant="ghost" onClick={() => handleDelete(wh.id)}><Trash2 className="w-4 h-4 text-red-500" /></Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        <Modal isOpen={modalOpen} onClose={() => setModalOpen(false)} title="Subscribe Webhook">
          <div className="space-y-4">
            <div><label className="block text-sm font-medium mb-1">Event</label><Select value={formEvent} onChange={(v) => setFormEvent(v as string)} options={eventOptions} /></div>
            <div><label className="block text-sm font-medium mb-1">URL (https://)</label><Input value={formUrl} onChange={(e) => setFormUrl(e.target.value)} placeholder="https://example.com/webhook" /></div>
            <div className="grid grid-cols-2 gap-2">
              <div><label className="block text-sm font-medium mb-1">Retry</label><Select value={formRetries} onChange={(v) => setFormRetries(v as string)} options={[{ value: '3' }, { value: '5' }, { value: '10' }]} /></div>
              <div><label className="block text-sm font-medium mb-1">Interval (s)</label><Select value={formInterval} onChange={(v) => setFormInterval(v as string)} options={[{ value: '30' }, { value: '60' }, { value: '120' }]} /></div>
            </div>
            <div className="flex justify-end gap-2 pt-2"><Button variant="ghost" onClick={() => setModalOpen(false)}>Batal</Button><Button onClick={handleCreate}>Subscribe</Button></div>
          </div>
        </Modal>
      </CardContent>
    </Card>
  );
}
