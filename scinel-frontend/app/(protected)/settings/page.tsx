'use client';
import { ThemeSelector } from '@/components/settings/ThemeSelector';
import { ApiKeyManager } from '@/components/settings/ApiKeyManager';
import { WebhookManager } from '@/components/settings/WebhookManager';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Checkbox } from '@/components/ui/Checkbox';
import { useState } from 'react';

export default function SettingsPage() {
  const [emailAlerts, setEmailAlerts] = useState(localStorage.getItem('notif_email_alerts') === 'true');
  const [lowStockAlerts, setLowStockAlerts] = useState(localStorage.getItem('notif_low_stock') === 'true');

  const handleNotifChange = (key: string, value: boolean) => {
    localStorage.setItem(key, String(value));
    if (key === 'notif_email_alerts') setEmailAlerts(value);
    if (key === 'notif_low_stock') setLowStockAlerts(value);
  };

  return (
    <div className="p-6 max-w-4xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">⚙️ Pengaturan</h1>
      <ThemeSelector />
      <ApiKeyManager />
      <WebhookManager />
      <Card>
        <CardHeader><CardTitle className="text-lg">🔔 Notifikasi</CardTitle></CardHeader>
        <CardContent className="space-y-2">
          <Checkbox checked={emailAlerts} onChange={(e) => handleNotifChange('notif_email_alerts', e.target.checked)} label="Email alert untuk anomali" />
          <Checkbox checked={lowStockAlerts} onChange={(e) => handleNotifChange('notif_low_stock', e.target.checked)} label="Email alert untuk stok menipis" />
          <p className="text-xs text-gray-400 mt-2">Preferensi disimpan di browser.</p>
        </CardContent>
      </Card>
    </div>
  );
}
