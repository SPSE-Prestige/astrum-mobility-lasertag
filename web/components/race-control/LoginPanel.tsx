"use client";

import { useState } from "react";
import { LogIn, LogOut } from "lucide-react";

interface LoginPanelProps {
  isAuthenticated: boolean;
  username: string | null;
  error: string | null;
  onLogin: (username: string, password: string) => boolean;
  onLogout: () => void;
}

export const LoginPanel = ({ isAuthenticated, username, error, onLogin, onLogout }: LoginPanelProps) => {
  const [form, setForm] = useState({ username: "", password: "" });

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const success = onLogin(form.username, form.password);
    if (success) {
      setForm({ username: "", password: "" });
    }
  };

  if (isAuthenticated) {
    return (
      <section className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
        <p className="text-xs uppercase tracking-[0.18em] text-zinc-500">Přihlášení</p>
        <p className="mt-2 text-sm text-zinc-200">
          Přihlášen: <span className="font-semibold text-[#00ff00]">{username}</span>
        </p>
        <button
          type="button"
          onClick={onLogout}
          className="mt-3 inline-flex items-center gap-2 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-xs font-semibold uppercase tracking-[0.15em] text-zinc-200 transition hover:border-zinc-500"
        >
          <LogOut className="h-4 w-4" />
          Odhlásit
        </button>
      </section>
    );
  }

  return (
    <section className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
      <p className="text-xs uppercase tracking-[0.18em] text-zinc-500">Přihlášení</p>
      <form className="mt-3 space-y-3" onSubmit={handleSubmit}>
        <label className="block text-xs uppercase tracking-[0.14em] text-zinc-500">
          Uživatelské jméno
          <input
            value={form.username}
            onChange={(event) => setForm((prev) => ({ ...prev, username: event.target.value }))}
            className="mt-1 w-full rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
          />
        </label>
        <label className="block text-xs uppercase tracking-[0.14em] text-zinc-500">
          Heslo
          <input
            type="password"
            value={form.password}
            onChange={(event) => setForm((prev) => ({ ...prev, password: event.target.value }))}
            className="mt-1 w-full rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
          />
        </label>
        {error && <p className="text-xs text-[#ff0000]">{error}</p>}
        <button
          type="submit"
          className="inline-flex items-center gap-2 rounded-md border border-[#00ff00]/70 bg-[#00ff00]/10 px-3 py-2 text-xs font-semibold uppercase tracking-[0.15em] text-[#00ff00] transition hover:bg-[#00ff00]/20"
        >
          <LogIn className="h-4 w-4" />
          Přihlásit
        </button>
      </form>
      <p className="mt-2 text-[11px] text-zinc-400">Přihlášení je povinné pro vstup do dashboardu.</p>
      <p className="mt-3 text-[11px] text-zinc-500">Demo účet: admin / admin123</p>
    </section>
  );
};