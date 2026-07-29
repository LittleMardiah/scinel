'use client';
import { useTheme } from '@/context/ThemeContext';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Radio } from '@/components/ui/Radio';

export function ThemeSelector() {
  const { theme, setTheme, isDark } = useTheme();
  const options = [
    { value: 'light', label: '☀️ Terang' },
    { value: 'dark', label: '🌙 Gelap' },
    { value: 'system', label: '🖥️ Sistem' },
  ];
  return (
    <Card>
      <CardHeader><CardTitle className="text-lg">Tema Tampilan</CardTitle></CardHeader>
      <CardContent className="space-y-2">
        {options.map((opt) => (
          <Radio key={opt.value} name="theme" value={opt.value} checked={theme === opt.value} onChange={() => setTheme(opt.value as 'light' | 'dark' | 'system')} label={opt.label} />
        ))}
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">{isDark ? '🌙 Mode gelap aktif' : '☀️ Mode terang aktif'}</p>
      </CardContent>
    </Card>
  );
}
