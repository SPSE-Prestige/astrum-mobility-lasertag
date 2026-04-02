import React from "react";
import { View, Text, StyleSheet, FlatList } from "react-native";
import Theme from "@/lib/theme";
import GlassCard from "@/components/ui/GlassCard";
import type { GameEventDTO } from "@/lib/api/types";
import { t, type Language } from "@/lib/i18n";

interface KillFeedProps {
  events: GameEventDTO[];
  currentNickname: string;
  lang: Language;
}

export default function KillFeed({ events, currentNickname, lang }: KillFeedProps) {
  const killEvents = events
    .filter((e) => e.type === "kill")
    .reverse()
    .slice(0, 15);

  return (
    <GlassCard style={styles.card} padding="md">
      <Text style={styles.title}>{t("killfeed.title", lang)}</Text>

      {killEvents.length === 0 ? (
        <Text style={styles.empty}>{t("killfeed.empty", lang)}</Text>
      ) : (
        <FlatList
          data={killEvents}
          scrollEnabled={false}
          keyExtractor={(e) => e.id}
          renderItem={({ item }) => {
            const attacker = String(item.payload.attacker_nickname ?? "???");
            const victim = String(item.payload.victim_nickname ?? "???");
            const isAttacker = attacker === currentNickname;
            const isVictim = victim === currentNickname;
            const time = new Date(item.timestamp).toLocaleTimeString("cs-CZ", {
              hour: "2-digit",
              minute: "2-digit",
              second: "2-digit",
            });

            return (
              <View style={styles.eventRow}>
                <Text style={styles.time}>{time}</Text>
                <Text
                  style={[
                    styles.attacker,
                    isAttacker && styles.highlight,
                  ]}
                  numberOfLines={1}
                >
                  {attacker}
                </Text>
                <Text style={styles.arrow}>→</Text>
                <Text
                  style={[
                    styles.victim,
                    isVictim && styles.highlightDanger,
                  ]}
                  numberOfLines={1}
                >
                  {victim}
                </Text>
              </View>
            );
          }}
        />
      )}
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
  empty: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textMuted,
    textAlign: "center",
    paddingVertical: Theme.spacing.lg,
  },
  eventRow: {
    flexDirection: "row",
    alignItems: "center",
    paddingVertical: Theme.spacing.xs,
    paddingHorizontal: Theme.spacing.sm,
    gap: Theme.spacing.sm,
  },
  time: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.xs,
    color: Theme.colors.textMuted,
    width: 58,
  },
  attacker: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.success,
    flex: 1,
  },
  arrow: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textMuted,
  },
  victim: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.danger,
    flex: 1,
    textAlign: "right",
  },
  highlight: {
    color: Theme.colors.brand,
    fontFamily: Theme.fontFamily.display,
  },
  highlightDanger: {
    color: Theme.colors.brand,
    fontFamily: Theme.fontFamily.display,
  },
});
