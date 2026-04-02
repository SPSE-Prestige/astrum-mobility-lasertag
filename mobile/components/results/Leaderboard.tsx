import React from "react";
import { View, Text, StyleSheet, FlatList } from "react-native";
import Theme from "@/lib/theme";
import GlassCard from "@/components/ui/GlassCard";
import type { LeaderboardPlayerDTO } from "@/lib/api/types";
import { t, type Language } from "@/lib/i18n";

interface LeaderboardProps {
  players: LeaderboardPlayerDTO[];
  lang: Language;
}

function getRankStyle(index: number) {
  if (index === 0) return { color: Theme.colors.brand, fontFamily: Theme.fontFamily.display };
  if (index === 1) return { color: Theme.colors.textSecondary };
  if (index === 2) return { color: "#cd7f32" }; // bronze
  return { color: Theme.colors.textMuted };
}

export default function Leaderboard({ players, lang }: LeaderboardProps) {
  if (!players || players.length === 0) return null;

  return (
    <GlassCard style={styles.card} padding="md">
      <Text style={styles.title}>{t("leaderboard.title", lang)}</Text>

      {/* Header */}
      <View style={styles.headerRow}>
        <Text style={[styles.headerCell, styles.rankCol]}>{t("leaderboard.rank", lang)}</Text>
        <Text style={[styles.headerCell, styles.nameCol]}>{t("leaderboard.player", lang)}</Text>
        <Text style={[styles.headerCell, styles.scoreCol]}>{t("leaderboard.score", lang)}</Text>
        <Text style={[styles.headerCell, styles.kdCol]}>{t("leaderboard.kd", lang)}</Text>
      </View>

      <View style={styles.divider} />

      <FlatList
        data={players}
        scrollEnabled={false}
        keyExtractor={(_, i) => i.toString()}
        renderItem={({ item, index }) => {
          const kd = item.deaths > 0 ? (item.kills / item.deaths).toFixed(1) : item.kills.toFixed(1);
          const rankStyle = getRankStyle(index);

          return (
            <View style={[styles.row, item.is_current && styles.currentRow]}>
              <Text style={[styles.cell, styles.rankCol, rankStyle]}>
                {index + 1}
              </Text>
              <View style={[styles.nameCol, styles.nameCell]}>
                <Text
                  style={[styles.cell, item.is_current && styles.currentName]}
                  numberOfLines={1}
                >
                  {item.nickname}
                </Text>
                {item.is_current && (
                  <Text style={styles.youBadge}>{t("common.you", lang)}</Text>
                )}
              </View>
              <Text style={[styles.cell, styles.scoreCol, { color: Theme.colors.textPrimary }]}>
                {item.score.toLocaleString()}
              </Text>
              <Text style={[styles.cell, styles.kdCol, { color: Theme.colors.textSecondary }]}>
                {kd}
              </Text>
            </View>
          );
        }}
      />
    </GlassCard>
  );
}

const styles = StyleSheet.create({
  card: {
    overflow: "hidden",
  },
  title: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textSecondary,
    letterSpacing: 2,
    textTransform: "uppercase",
    marginBottom: Theme.spacing.md,
    paddingHorizontal: Theme.spacing.sm,
  },
  headerRow: {
    flexDirection: "row",
    paddingHorizontal: Theme.spacing.sm,
    paddingBottom: Theme.spacing.xs,
  },
  headerCell: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xs,
    color: Theme.colors.textMuted,
    letterSpacing: 1,
    textTransform: "uppercase",
  },
  divider: {
    height: 1,
    backgroundColor: Theme.colors.border,
    marginBottom: Theme.spacing.xs,
  },
  row: {
    flexDirection: "row",
    alignItems: "center",
    paddingVertical: Theme.spacing.sm,
    paddingHorizontal: Theme.spacing.sm,
    borderRadius: Theme.borderRadius.sm,
  },
  currentRow: {
    backgroundColor: Theme.colors.brandMuted,
  },
  cell: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.md,
    color: Theme.colors.textSecondary,
  },
  rankCol: {
    width: 32,
    textAlign: "center",
  },
  nameCol: {
    flex: 1,
    flexDirection: "row",
    alignItems: "center",
    gap: Theme.spacing.xs,
  },
  nameCell: {
    flexDirection: "row",
    alignItems: "center",
    gap: Theme.spacing.xs,
  },
  currentName: {
    color: Theme.colors.brand,
    fontFamily: Theme.fontFamily.display,
  },
  youBadge: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xs,
    color: Theme.colors.brand,
    backgroundColor: Theme.colors.brandMuted,
    paddingHorizontal: 6,
    paddingVertical: 1,
    borderRadius: 4,
    letterSpacing: 1,
    overflow: "hidden",
  },
  scoreCol: {
    width: 70,
    textAlign: "right",
  },
  kdCol: {
    width: 48,
    textAlign: "right",
  },
});
