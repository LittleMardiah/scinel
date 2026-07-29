'use client';

import { useEffect, useState } from 'react';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface ToastMessage {
  id: number;
  type: ToastType;
  message: string;
}

let toastId = 0;
const listeners: ((messages: ToastMessage[]) => void)[] = [];
let messages: ToastMessage[] = [];

const notifyListeners = () => {
  listeners.forEach((listener) => listener([...messages]));
};

export const Toast = {
  success: (message: string) => {
    messages.push({ id: ++toastId, type: 'success', message });
    notifyListeners();
    setTimeout(() => {
      messages = messages.filter((m) => m.id !== toastId);
      notifyListeners();
    }, 5000);
  },
  error: (message: string) => {
    messages.push({ id: ++toastId, type: 'error', message });
    notifyListeners();
    setTimeout(() => {
      messages = messages.filter((m) => m.id !== toastId);
      notifyListeners();
    }, 5000);
  },
  warning: (message: string) => {
    messages.push({ id: ++toastId, type: 'warning', message });
    notifyListeners();
    setTimeout(() => {
      messages = messages.filter((m) => m.id !== toastId);
      notifyListeners();
    }, 5000);
  },
  info: (message: string) => {
    messages.push({ id: ++toastId, type: 'info', message });
    notifyListeners();
    setTimeout(() => {
      messages = messages.filter((m) => m.id !== toastId);
      notifyListeners();
    }, 5000);
  },
  subscribe: (listener: (messages: ToastMessage[]) => void) => {
    listeners.push(listener);
    return () => {
      const index = listeners.indexOf(listener);
      if (index > -1) listeners.splice(index, 1);
    };
  },
  getMessages: () => messages,
};

export default Toast;
