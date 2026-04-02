import { useEffect, useState, useRef, useCallback } from 'react';
import {
  SafeAreaView,
  Text,
  View,
  TouchableOpacity,
  ActivityIndicator,
  ScrollView,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { useGameSession } from '../lib/useGameSession';

export default function DashboardScreen() {
  const { code } = useLocalSearchParams<{ code: string }>();
  const router = useRouter();
  const { session, loading, error } = useGameSession(code ?? '');
  const [displayTime, setDisplayTime] = useState<number | null>(null);
  const lastServerTime = useRef<number | null>(null);

  // Sync timer with server data, then count down locally
  useEffect(() => {
    if (session?.remaining_time != null) {
      lastServerTime.current = session.remaining_time;
      setDisplayTime(session.remaining_time);
    }
  }, [session?.remaining_time]);

  useEffect(() => {
    if (displayTime == null || displayTime <= 0) return;
    const timer = setInterval(() => {
      setDisplayTime((prev) => (prev != null && prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => clearInterval(timer);
  }, [displayTime != null && displayTime > 0]);

  // Auto-navigate to results when game finishes
  useEffect(() => {
    if (session?.game.status === 'finished') {
      router.replace({ pathname: '/results', params: { code: code ?? '' } });
    }
  }, [session?.game.status]);

  const disconnect = useCallback(() => {
    router.replace('/');
  }, [router]);

  if (loading && !session) {
    return (
      <SafeAreaView className="flex-1 bg-[#020303] justify-center items-center">
        <ActivityIndicator size="large" color="#00F0FF" />
        <Text className="text-[#00F0FF] mt-4">Načítání...</Text>
      </SafeAreaView>
    );
  }

  if (error && !session) {
    return (
      <SafeAreaView className="flex-1 bg-[#020303] justify-center items-center px-6">
        <Text className="text-[#FF3B30] text-lg mb-4">{error}</Text>
        <TouchableOpacity
          onPress={disconnect}
          className="bg-[#FF3B30]/20 px-6 py-3 rounded-xl"
        >
          <Text className="text-[#FF3B30] font-bold">ZPĚT</Text>
        </TouchableOpacity>
      </SafeAreaView>
    );
  }

  if (!session) return null;

  const { player, game, team } = session;
  const accuracy =
    player.shots_fired > 0
      ? Math.round((player.kills / player.shots_fired) * 100)
      : 0;

  const formatTime = (seconds: number) => {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}:${s.toString().padStart(2, '0')}`;
  };

  return (
    <SafeAreaView className="flex-1 bg-[#020303]">
      <ScrollView className="flex-1" contentContainerClassName="pb-6">
        {/* Top Bar */}
        <View className="flex-row items-center justify-between px-4 py-3 border-b border-[#1a1a2e]">
          <View className="flex-row items-center flex-1">
            <View
              className="w-3 h-3 rounded-full mr-2"
              style={{ backgroundColor: team.color }}
            />
            <Text
              className="text-white text-lg font-bold mr-2"
              numberOfLines={1}
            >
              {player.nickname}
            </Text>
            <Text className="text-gray-500 text-xs">{team.name}</Text>
          </View>
          <View className="bg-[#0A0F0F] px-3 py-1 rounded-lg border border-[#1a1a2e]">
            <Text className="text-[#00F0FF] text-xs font-mono">
              {game.code}
            </Text>
          </View>
        </View>

        {/* Game Timer */}
        <View className="items-center py-4">
          <Text className="text-gray-500 text-xs uppercase tracking-widest mb-1">
            Zbývající čas
          </Text>
          <Text className="text-white text-4xl font-bold font-mono">
            {displayTime != null ? formatTime(displayTime) : '--:--'}
          </Text>
        </View>

        {/* Alive / Dead Status */}
        <View className="items-center mb-4">
          {player.is_alive ? (
            <View className="bg-[#30D158]/15 px-6 py-2 rounded-full border border-[#30D158]/30">
              <Text className="text-[#30D158] font-bold text-sm tracking-widest">
                ● ALIVE
              </Text>
            </View>
          ) : (
            <View className="bg-[#FF3B30]/15 px-6 py-2 rounded-full border border-[#FF3B30]/30">
              <Text className="text-[#FF3B30] font-bold text-sm tracking-widest">
                ✕ DEAD — respawn {game.settings.respawn_delay}s
              </Text>
            </View>
          )}
        </View>

        {/* Stats Grid 2x2 */}
        <View className="px-4 mb-4">
          <View className="flex-row mb-3">
            {/* Score */}
            <View className="flex-1 bg-[#0A0F0F] rounded-xl border border-[#1a1a2e] p-4 mr-1.5">
              <Text className="text-gray-500 text-xs uppercase tracking-wider mb-1">
                Skóre
              </Text>
              <Text className="text-[#00F0FF] text-3xl font-bold">
                {player.score}
              </Text>
            </View>
            {/* K/D */}
            <View className="flex-1 bg-[#0A0F0F] rounded-xl border border-[#1a1a2e] p-4 ml-1.5">
              <Text className="text-gray-500 text-xs uppercase tracking-wider mb-1">
                Kills / Deaths
              </Text>
              <Text className="text-white text-3xl font-bold">
                {player.kills}
                <Text className="text-gray-500"> / </Text>
                {player.deaths}
              </Text>
            </View>
          </View>
          <View className="flex-row">
            {/* Kill Streak */}
            <View className="flex-1 bg-[#0A0F0F] rounded-xl border border-[#1a1a2e] p-4 mr-1.5">
              <Text className="text-gray-500 text-xs uppercase tracking-wider mb-1">
                🔥 Kill Streak
              </Text>
              <Text className="text-orange-400 text-3xl font-bold">
                {player.kill_streak}
              </Text>
            </View>
            {/* Weapon Level */}
            <View className="flex-1 bg-[#0A0F0F] rounded-xl border border-[#1a1a2e] p-4 ml-1.5">
              <Text className="text-gray-500 text-xs uppercase tracking-wider mb-1">
                ⚡ Weapon Level
              </Text>
              <Text className="text-yellow-400 text-3xl font-bold">
                {player.weapon_level}
              </Text>
            </View>
          </View>
        </View>

        {/* Accuracy Bar */}
        <View className="px-4 mb-4">
          <View className="bg-[#0A0F0F] rounded-xl border border-[#1a1a2e] p-4">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-gray-500 text-xs uppercase tracking-wider">
                🎯 Přesnost
              </Text>
              <Text className="text-white font-bold">{accuracy}%</Text>
            </View>
            <View className="h-2 bg-[#1a1a2e] rounded-full overflow-hidden">
              <View
                className="h-full rounded-full bg-[#00F0FF]"
                style={{ width: `${Math.min(accuracy, 100)}%` }}
              />
            </View>
            <Text className="text-gray-600 text-xs mt-1">
              {player.kills} zásahů / {player.shots_fired} výstřelů
            </Text>
          </View>
        </View>

        {/* Team Info */}
        <View className="px-4 mb-6">
          <View className="bg-[#0A0F0F] rounded-xl border border-[#1a1a2e] p-4 flex-row items-center">
            <View
              className="w-5 h-5 rounded-full mr-3"
              style={{ backgroundColor: team.color }}
            />
            <View>
              <Text className="text-gray-500 text-xs uppercase tracking-wider">
                Tým
              </Text>
              <Text className="text-white font-bold text-base">
                {team.name}
              </Text>
            </View>
          </View>
        </View>

        {/* Disconnect */}
        <View className="px-4">
          <TouchableOpacity
            onPress={disconnect}
            className="bg-[#FF3B30]/10 border border-[#FF3B30]/30 py-3 rounded-xl items-center"
          >
            <Text className="text-[#FF3B30] font-bold tracking-wider">
              ODPOJIT SE
            </Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
