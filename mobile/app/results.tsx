import { useEffect, useState } from 'react';
import {
  SafeAreaView,
  Text,
  View,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { getPlayerSession } from '../lib/api';
import type { PlayerSession } from '../lib/types';

export default function ResultsScreen() {
  const { code } = useLocalSearchParams<{ code: string }>();
  const router = useRouter();
  const [session, setSession] = useState<PlayerSession | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!code) return;
    getPlayerSession(code)
      .then(setSession)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [code]);

  if (loading) {
    return (
      <SafeAreaView className="flex-1 bg-[#020303] justify-center items-center">
        <ActivityIndicator size="large" color="#00F0FF" />
      </SafeAreaView>
    );
  }

  if (!session) {
    return (
      <SafeAreaView className="flex-1 bg-[#020303] justify-center items-center px-6">
        <Text className="text-[#FF3B30] text-lg mb-4">
          Nepodařilo se načíst výsledky
        </Text>
        <TouchableOpacity
          onPress={() => router.replace('/')}
          className="bg-[#00F0FF] px-6 py-3 rounded-xl"
        >
          <Text className="text-[#020303] font-bold">ZPĚT</Text>
        </TouchableOpacity>
      </SafeAreaView>
    );
  }

  const { player, team } = session;
  const accuracy =
    player.shots_fired > 0
      ? Math.round((player.kills / player.shots_fired) * 100)
      : 0;
  const kd =
    player.deaths > 0 ? (player.kills / player.deaths).toFixed(2) : player.kills.toFixed(2);

  return (
    <SafeAreaView className="flex-1 bg-[#020303]">
      <View className="flex-1 justify-center px-6">
        {/* Header */}
        <View className="items-center mb-8">
          <Text className="text-5xl mb-2">🏁</Text>
          <Text className="text-white text-3xl font-bold tracking-widest">
            HRA SKONČILA
          </Text>
          <View className="flex-row items-center mt-2">
            <View
              className="w-3 h-3 rounded-full mr-2"
              style={{ backgroundColor: team.color }}
            />
            <Text className="text-gray-400">
              {player.nickname} • {team.name}
            </Text>
          </View>
        </View>

        {/* Final Score */}
        <View className="bg-[#0A0F0F] rounded-2xl border border-[#00F0FF]/30 p-6 mb-4 items-center">
          <Text className="text-gray-500 text-xs uppercase tracking-widest mb-1">
            Celkové skóre
          </Text>
          <Text className="text-[#00F0FF] text-5xl font-bold">
            {player.score}
          </Text>
        </View>

        {/* Stats Grid */}
        <View className="flex-row mb-3">
          <StatCard label="Kills" value={String(player.kills)} />
          <View className="w-3" />
          <StatCard label="Deaths" value={String(player.deaths)} />
        </View>
        <View className="flex-row mb-3">
          <StatCard label="K/D Ratio" value={kd} />
          <View className="w-3" />
          <StatCard label="Přesnost" value={`${accuracy}%`} />
        </View>
        <View className="flex-row mb-8">
          <StatCard label="⚡ Weapon Lvl" value={String(player.weapon_level)} />
          <View className="w-3" />
          <StatCard label="🔥 Best Streak" value={String(player.kill_streak)} />
        </View>

        {/* New Game Button */}
        <TouchableOpacity
          onPress={() => router.replace('/')}
          className="bg-[#00F0FF] py-4 rounded-xl items-center"
        >
          <Text className="text-[#020303] text-lg font-bold tracking-wider">
            NOVÁ HRA
          </Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}

function StatCard({ label, value }: { label: string; value: string }) {
  return (
    <View className="flex-1 bg-[#0A0F0F] rounded-xl border border-[#1a1a2e] p-4 items-center">
      <Text className="text-gray-500 text-xs uppercase tracking-wider mb-1">
        {label}
      </Text>
      <Text className="text-white text-2xl font-bold">{value}</Text>
    </View>
  );
}
