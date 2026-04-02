import { Wifi, WifiOff } from "lucide-react";

interface ConnectionStatusProps {
  status: "connecting" | "connected" | "disconnected";
}

const config = {
  connecting: { icon: Wifi, label: "Connecting...", dot: "bg-amber-500 animate-pulse", text: "text-amber-400" },
  connected: { icon: Wifi, label: "Live", dot: "bg-red-500", text: "text-red-400" },
  disconnected: { icon: WifiOff, label: "Offline", dot: "bg-red-300", text: "text-red-300" },
};

export function ConnectionStatus({ status }: ConnectionStatusProps) {
  const { icon: Icon, label, dot, text } = config[status];

  return (
    <div className={`inline-flex items-center gap-1.5 text-xs ${text}`} role="status" aria-label={`Connection: ${label}`}>
      <span className={`h-1.5 w-1.5 rounded-full ${dot}`} />
      <Icon className="h-3 w-3" />
      <span className="uppercase tracking-wider">{label}</span>
    </div>
  );
}
