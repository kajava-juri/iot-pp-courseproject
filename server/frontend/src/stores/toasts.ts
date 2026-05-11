import { writable } from 'svelte/store';

type Toast = {
  id: string;
  message: string;
  type?: 'success' | 'info' | 'error';
};

const toasts = writable<Toast[]>([]);

function addToast(message: string, type: Toast['type'] = 'success', duration = 3500) {
  const id = crypto.randomUUID();
  const toast: Toast = { id, message, type };
  toasts.update((t) => [...t, toast]);
  if (duration > 0) {
    setTimeout(() => {
      removeToast(id);
    }, duration);
  }
  return id;
}

function removeToast(id: string) {
  toasts.update((t) => t.filter((x) => x.id !== id));
}

export { toasts, addToast, removeToast };
