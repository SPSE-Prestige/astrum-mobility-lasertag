import { API_BASE } from './config';
import type { PlayerSession } from './types';

export async function getPlayerSession(code: string): Promise<PlayerSession> {
  const res = await fetch(`${API_BASE}/api/player/session/${code}`);
  if (!res.ok) throw new Error('Invalid session code');
  return res.json();
}
