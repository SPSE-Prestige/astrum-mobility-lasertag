import React from "react";
import { View, Text, StyleSheet } from "react-native";
import Theme from "@/lib/theme";
import GlassCard from "@/components/ui/GlassCard";
import type { PlayerDTO } from "@/lib/api/types";
import { t, type Language } from "@/lib/i18n";

interface StatsGridProps {
  player: PlayerDTO;
  lang: Language;
}

interface StatItem {
  label: string;
  value: string;
  color: string;
}

export default function StatsGrid({ player, lang }: StatsGridProps) {
  const accuracy =
    player.shots_fired > 0
      ? ((player.kills / player.shots_fired) * 100).toFixed(1)
      : "0.0";

  const kdRatio =
    player.deaths > 0 ? (player.kills / player.deaths).toFixed(2) : player.kills.toFixed(2);

  const stats: StatItem[] = [
    {
      label: t("stats.kills", lang),
      value: player.kills.toString(),
      color: Theme.colors.success,
    },
    {
      label: t("stats.deaths", lang),
      value: player.deaths.toString(),
      color: Theme.colors.danger,
    },
    {
      label: t("stats.accuracy", lang),
      value: `${accuracy}%`,
      color: Theme.colors.info,
    },
    {
      label: t("stats.kd", lang),
      value: kdRatio,
      color: Theme.colors.warning,
    },
    {
      label: t("stats.shotsFired", lang),
      value: player.shots_fired.toLocaleString(),
      color: Theme.colors.textSecondary,
    },
    {
      label: t("stats.bestStreak", lang),
      value: player.kill_streak.toString(),
      color: Theme.colors.warning,
    },
  ];

  return (
    <View style={styles.grid}>
      {stats.map((stat) => (
        <GlassCard key={stat.label} style={styles.cell} padding="md">
          <Text style={[styles.value, { color: stat.color }]}>{stat.value}</Text>
          <Text style={styles.label}>{stat.label}</Text>
        </GlassCard>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  grid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: Theme.spacing.sm,
  },
  cell: {
    width: "31.5%",
    flexGrow: 1,
    alignItems: "center",
    minWidth: 100,
  },
  value: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xxl,
    letterSpacing: 1,
    lineHeight: Theme.fontSize.xxl * 1.2,
  },
  label: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.xs,
    color: Theme.colors.textMuted,
    letterSpacing: 1,
    textTransform: "uppercase",
    marginTop: Theme.spacing.xxs,
  },
});
