'use client';

import React, { useState } from 'react';
import { Modal } from '@/components/common/Modal';
import { Button } from '@/components/common/Button';
import { TextArea } from '@/components/ui/TextArea';
import { Input } from '@/components/common/Input';

interface ResolutionModalProps {
  isOpen: boolean;
  onClose: () => void;
  onResolve: (notes: string) => Promise<void>;
  anomalyId: string;
  isLoading?: boolean;
  currentStatus?: string;
}

export const ResolutionModal: React.FC<ResolutionModalProps> = ({
  isOpen,
  onClose,
  onResolve,
  anomalyId,
  isLoading,
  currentStatus,
}) => {
  const [notes, setNotes] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (!notes.trim()) {
      setError('Catatan resolusi diperlukan');
      return;
    }
    try {
      await onResolve(notes);
      setNotes('');
      onClose();
    } catch (err: any) {
      setError(err.message || 'Gagal menyelesaikan anomali');
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Selesaikan Anomali">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <p className="text-sm text-gray-600 dark:text-gray-300">
            Anomali ID: <span className="font-mono text-xs">{anomalyId}</span>
          </p>
          {currentStatus && (
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Status saat ini: <span className="font-medium">{currentStatus}</span>
            </p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Catatan Resolusi <span className="text-red-500">*</span>
          </label>
          <TextArea
            placeholder="Jelaskan bagaimana anomali ini diselesaikan..."
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            rows={4}
            error={error || undefined}
            required
          />
          <p className="mt-1 text-xs text-gray-400">
            {notes.length}/500 karakter
          </p>
        </div>

        <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
          <Button type="button" variant="ghost" onClick={onClose} disabled={isLoading}>
            Batal
          </Button>
          <Button type="submit" isLoading={isLoading} disabled={isLoading || !notes.trim()}>
            Selesaikan
          </Button>
        </div>
      </form>
    </Modal>
  );
};
