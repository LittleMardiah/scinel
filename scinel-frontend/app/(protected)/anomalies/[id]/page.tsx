'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { SeverityBadge } from '@/components/anomaly/SeverityBadge';
import { AnomalyTimeline } from '@/components/anomaly/AnomalyTimeline';
import { ResolutionModal } from '@/components/anomaly/ResolutionModal';
import { Breadcrumb } from '@/components/common/Breadcrumb';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Skeleton } from '@/components/common/Skeleton';
import { Badge } from '@/components/ui/Badge';
import { ArrowLeft, AlertCircle, CheckCircle, Clock } from 'lucide-react';
import { Toast } from '@/components/common/Toast';
import { apiClient } from '@/lib/api';

export default function AnomalyDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuth();
  const anomalyId = params?.id as string;

  const [anomaly, setAnomaly] = useState<any>(null);
  const [timeline, setTimeline] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showResolutionModal, setShowResolutionModal] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  const fetchData = async () => {
    if (!anomalyId) return;
    setIsLoading(true);
    setError(null);
    try {
      const [detailRes, timelineRes] = await Promise.all([
        apiClient.get(`/anomalies/${anomalyId}`),
        apiClient.get(`/anomalies/${anomalyId}/history`),
      ]);
      setAnomaly(detailRes.data.data);
      setTimeline(timelineRes.data.data?.timeline || []);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat detail anomali');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [anomalyId]);

  const handleAcknowledge = async () => {
    if (!anomaly) return;
    setIsUpdating(true);
    try {
      await apiClient.put(`/anomalies/${anomalyId}`, { status: 'acknowledged' });
      Toast.success('Anomali sedang diteliti');
      fetchData();
    } catch (err: any) {
      Toast.error(err.message || 'Gagal mengupdate status');
    } finally {
      setIsUpdating(false);
    }
  };

  const handleResolve = async (notes: string) => {
    if (!anomaly) return;
    setIsUpdating(true);
    try {
      await apiClient.put(`/anomalies/${anomalyId}`, {
        status: 'resolved',
        resolution_notes: notes,
      });
      Toast.success('Anomali selesai');
      fetchData();
      setShowResolutionModal(false);
    } catch (err: any) {
      Toast.error(err.message || 'Gagal mengupdate status');
      throw err;
    } finally {
      setIsUpdating(false);
    }
  };

  const breadcrumbItems = [
    { label: 'Dashboard', href: '/dashboard' },
    { label: 'Anomalies', href: '/anomalies' },
    { label: `Detail: ${anomaly?.id?.slice(0, 8) || '...'}` },
  ];

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-48 w-full" />
          <Skeleton className="h-48 w-full" />
        </div>
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (error || !anomaly) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="rounded-lg bg-red-50 p-6 text-center dark:bg-red-900/20">
          <AlertCircle className="mx-auto h-12 w-12 text-red-600 dark:text-red-400" />
          <p className="mt-2 text-red-600 dark:text-red-400">{error || 'Anomali tidak ditemukan'}</p>
          <Button variant="primary" onClick={() => router.push('/anomalies')} className="mt-4">
            Kembali ke Daftar Anomali
          </Button>
        </div>
      </div>
    );
  }

  const statusColors = {
    open: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
    acknowledged: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
    resolved: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  };

  const statusLabels = {
    open: 'Terbuka',
    acknowledged: 'Diteliti',
    resolved: 'Selesai',
  };

  const canAcknowledge = anomaly.status === 'open';
  const canResolve = anomaly.status === 'open' || anomaly.status === 'acknowledged';

  return (
    <div className="space-y-6">
      <Breadcrumb items={breadcrumbItems} />

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/anomalies">
            <Button variant="ghost" size="sm" className="gap-2">
              <ArrowLeft className="h-4 w-4" />
              Kembali
            </Button>
          </Link>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Detail Anomali
          </h1>
        </div>
        <div className="flex gap-2">
          {canAcknowledge && (
            <Button
              variant="secondary"
              size="sm"
              onClick={handleAcknowledge}
              disabled={isUpdating}
              className="gap-2"
            >
              <Clock className="h-4 w-4" />
              {isUpdating ? 'Memproses...' : 'Teliti'}
            </Button>
          )}
          {canResolve && (
            <Button
              variant="primary"
              size="sm"
              onClick={() => setShowResolutionModal(true)}
              disabled={isUpdating}
              className="gap-2"
            >
              <CheckCircle className="h-4 w-4" />
              Selesaikan
            </Button>
          )}
        </div>
      </div>

      {/* Detail Cards */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Informasi Anomali</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-500 dark:text-gray-400">ID</span>
              <span className="font-mono text-sm text-gray-900 dark:text-white">{anomaly.id}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500 dark:text-gray-400">SKU</span>
              <span className="font-medium text-gray-900 dark:text-white">{anomaly.sku_code || '-'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500 dark:text-gray-400">Tipe</span>
              <span className="text-gray-900 dark:text-white">{anomaly.type || '-'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500 dark:text-gray-400">Severity</span>
              <SeverityBadge severity={anomaly.severity} />
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500 dark:text-gray-400">Status</span>
              <Badge className={statusColors[anomaly.status as keyof typeof statusColors]}>
                {statusLabels[anomaly.status as keyof typeof statusLabels] || anomaly.status}
              </Badge>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500 dark:text-gray-400">Z-Score</span>
              <span className="font-medium text-gray-900 dark:text-white">{anomaly.z_score?.toFixed(2) || '-'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500 dark:text-gray-400">Detected At</span>
              <span className="text-gray-900 dark:text-white">
                {new Date(anomaly.detected_at).toLocaleString('id-ID')}
              </span>
            </div>
            {anomaly.resolved_at && (
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Resolved At</span>
                <span className="text-gray-900 dark:text-white">
                  {new Date(anomaly.resolved_at).toLocaleString('id-ID')}
                </span>
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Deskripsi & Metrik</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <span className="text-sm text-gray-500 dark:text-gray-400">Deskripsi</span>
              <p className="mt-1 text-gray-900 dark:text-white">{anomaly.description || 'Tidak ada deskripsi'}</p>
            </div>
            <div className="grid grid-cols-2 gap-3 pt-3 border-t border-gray-200 dark:border-gray-700">
              <div>
                <span className="text-xs text-gray-500 dark:text-gray-400">Historical Average</span>
                <p className="text-sm font-medium text-gray-900 dark:text-white">
                  {anomaly.historical_average?.toFixed(2) || '-'}
                </p>
              </div>
              <div>
                <span className="text-xs text-gray-500 dark:text-gray-400">Detected Value</span>
                <p className="text-sm font-medium text-gray-900 dark:text-white">
                  {anomaly.detected_value?.toFixed(2) || '-'}
                </p>
              </div>
            </div>
            {anomaly.resolution_notes && (
              <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                <span className="text-sm text-gray-500 dark:text-gray-400">Resolution Notes</span>
                <p className="mt-1 text-gray-900 dark:text-white">{anomaly.resolution_notes}</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Timeline */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Timeline</CardTitle>
        </CardHeader>
        <CardContent>
          <AnomalyTimeline events={timeline} />
        </CardContent>
      </Card>

      {/* Resolution Modal */}
      <ResolutionModal
        isOpen={showResolutionModal}
        onClose={() => setShowResolutionModal(false)}
        onResolve={handleResolve}
        anomalyId={anomaly.id}
        isLoading={isUpdating}
        currentStatus={statusLabels[anomaly.status as keyof typeof statusLabels]}
      />
    </div>
  );
}
