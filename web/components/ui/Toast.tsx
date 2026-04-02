"use client";

import { useCallback, useEffect, useState } from "react";
import { AlertTriangle, CheckCircle, Info, X } from "lucide-react";

export type ToastType = "success" | "error" | "info";

interface Toast {
  id: string;
  type: ToastType;
  message: string;
}

let toastId = 0;
const listeners: Set<(toast: Toast) => void> = new Set();

export function showToast(type: ToastType, message: string) {
  const toast: Toast = { id: String(++toastId), type, message };
  listeners.forEach((fn) => fn(toast));
}

const icons: Record<ToastType, typeof AlertTriangle> = {
  success: CheckCircle,
  error: AlertTriangle,
  info: Info,
};

const styles: Record<ToastType, string> = {
  success: "border-emerald-500/50 bg-emerald-950/90 text-emerald-200",
  error: "border-red-500/50 bg-red-950/90 text-red-200",
  info: "border-blue-500/50 bg-blue-950/90 text-blue-200",
};

export function ToastContainer() {
  const [toasts, setToasts] = useState<Toast[]>([]);

  useEffect(() => {
    const handler = (toast: Toast) => {
      setToasts((prev) => [...prev.slice(-4), toast]);
    };
    listeners.add(handler);
    return () => { listeners.delete(handler); };
  }, []);

  // Auto-dismiss after 5s
  useEffect(() => {
    if (toasts.length === 0) return;
    const timer = setTimeout(() => {
      setToasts((prev) => prev.slice(1));
    }, 5000);
    return () => clearTimeout(timer);
  }, [toasts]);

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2" role="status" aria-live="polite">
      {toasts.map((toast) => {
        const Icon = icons[toast.type];
        return (
          <div
            key={toast.id}
            className={`flex items-center gap-3 rounded-lg border px-4 py-3 shadow-lg backdrop-blur animate-in slide-in-from-right ${styles[toast.type]}`}
          >
            <Icon className="h-4 w-4 shrink-0" />
            <span className="text-sm">{toast.message}</span>
            <button
              type="button"
              onClick={() => dismiss(toast.id)}
              className="ml-2 shrink-0 rounded p-1 opacity-60 hover:opacity-100 transition"
              aria-label="Dismiss"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        );
      })}
    </div>
  );
}
