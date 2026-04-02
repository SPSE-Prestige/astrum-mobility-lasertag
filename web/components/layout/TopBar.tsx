import { LogOut } from "lucide-react";
import { ConnectionStatus } from "@/components/ui/ConnectionStatus";
import { t } from "@/lib/i18n";
import type { Language } from "@/types/i18n";

interface TopBarProps {
  language: Language;
  username: string | null;
  wsStatus: "connecting" | "connected" | "disconnected";
  showConnectionStatus: boolean;
  onLanguageChange: (lang: Language) => void;
  onLogout: () => void;
  className?: string;
}

export function TopBar({ language, username, wsStatus, showConnectionStatus, onLanguageChange, onLogout, className = "" }: TopBarProps) {
  return (
    <div className={`flex items-center justify-between gap-3 rounded-xl border border-zinc-800 bg-black/30 px-3 py-2 ${className}`}>
      <div className="flex items-center gap-3">
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => onLanguageChange("cs")}
            className={`rounded-md border px-2 py-1 text-xs ${language === "cs" ? "border-[#ff0a0a] text-[#ff0a0a]" : "border-zinc-700 text-zinc-400"}`}
          >
            CZ
          </button>
          <button
            type="button"
            onClick={() => onLanguageChange("en")}
            className={`rounded-md border px-2 py-1 text-xs ${language === "en" ? "border-[#ff0a0a] text-[#ff0a0a]" : "border-zinc-700 text-zinc-400"}`}
          >
            EN
          </button>
        </div>
        {showConnectionStatus && <ConnectionStatus status={wsStatus} />}
      </div>
      <span className="text-xs uppercase tracking-[0.16em] text-zinc-400">{username ?? ""}</span>
      <button
        type="button"
        onClick={onLogout}
        className="inline-flex items-center gap-2 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.14em] text-zinc-200 transition hover:border-zinc-500"
      >
        <LogOut className="h-3.5 w-3.5" />
        {t("auth.logout", language)}
      </button>
    </div>
  );
}
