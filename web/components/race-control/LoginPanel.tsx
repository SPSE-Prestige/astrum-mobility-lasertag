"use client";

import { useState } from "react";
import { LogIn, LogOut } from "lucide-react";
import type { Language } from "@/types/i18n";

interface LoginPanelProps {
  isAuthenticated: boolean;
  username: string | null;
  error: string | null;
  language: Language;
  onLogin: (username: string, password: string) => boolean;
  onLogout: () => void;
}

export const LoginPanel = ({ isAuthenticated, username, error, language, onLogin, onLogout }: LoginPanelProps) => {
  const [form, setForm] = useState({ username: "", password: "" });
  const invalidCredentialsMessage = language === "cs" ? "Neplatné uživatelské jméno nebo heslo." : "Invalid username or password.";

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
        <p className="text-xs uppercase tracking-[0.18em] text-zinc-500">{language === "cs" ? "Přihlášení" : "Login"}</p>
        <p className="mt-2 text-sm text-zinc-200">
          {language === "cs" ? "Přihlášen" : "Logged in"}: <span className="font-semibold text-[#00ff00]">{username}</span>
        </p>
        <button
          type="button"
          onClick={onLogout}
          className="mt-3 inline-flex items-center gap-2 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-xs font-semibold uppercase tracking-[0.15em] text-zinc-200 transition hover:border-zinc-500"
        >
          <LogOut className="h-4 w-4" />
          {language === "cs" ? "Odhlásit" : "Logout"}
        </button>
      </section>
    );
  }

  return (
    <section className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
      <p className="text-xs uppercase tracking-[0.18em] text-zinc-500">{language === "cs" ? "Přihlášení" : "Login"}</p>
      <form className="mt-3 space-y-3" onSubmit={handleSubmit}>
        <label className="block text-xs uppercase tracking-[0.14em] text-zinc-500">
          {language === "cs" ? "Uživatelské jméno" : "Username"}
          <input
            value={form.username}
            onChange={(event) => setForm((prev) => ({ ...prev, username: event.target.value }))}
            className="mt-1 w-full rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
          />
        </label>
        <label className="block text-xs uppercase tracking-[0.14em] text-zinc-500">
          {language === "cs" ? "Heslo" : "Password"}
          <input
            type="password"
            value={form.password}
            onChange={(event) => setForm((prev) => ({ ...prev, password: event.target.value }))}
            className="mt-1 w-full rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
          />
        </label>
        {error && <p className="text-xs text-[#ff0000]">{error === "invalid_credentials" ? invalidCredentialsMessage : error}</p>}
        <button
          type="submit"
          className="inline-flex items-center gap-2 rounded-md border border-[#00ff00]/70 bg-[#00ff00]/10 px-3 py-2 text-xs font-semibold uppercase tracking-[0.15em] text-[#00ff00] transition hover:bg-[#00ff00]/20"
        >
          <LogIn className="h-4 w-4" />
          {language === "cs" ? "Přihlásit" : "Login"}
        </button>
      </form>
      <p className="mt-2 text-[11px] text-zinc-400">
        {language === "cs" ? "Přihlášení je povinné pro vstup do dashboardu." : "Login is required to access the dashboard."}
      </p>
      <p className="mt-3 text-[11px] text-zinc-500">{language === "cs" ? "Demo účet" : "Demo account"}: admin / admin123</p>
    </section>
  );
};