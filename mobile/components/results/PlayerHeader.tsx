import React from "react";
import { View, Text, StyleSheet } from "react-native";
import Theme from "@/lib/theme";
import Badge from "@/components/ui/Badge";
import type { PlayerDTO, GameDTO, TeamDTO } from "@/lib/api/types";
import { t, type Language } from "@/lib/i18n";

interface PlayerHeaderProps {
  player: PlayerDTO;
  game: GameDTO;
  team: TeamDTO | null;
  lang: Language;
}

export default function PlayerHeader({ player, game, team, lang }: PlayerHeaderProps) {
  const statusVariant =
    game.status === "running" ? "success" : game.status === "finished" ? "brand" : "muted";
  const statusLabel = t(
    `results.status.${game.status}` as `results.status.${"running" | "finished" | "lobby"}`,
    lang
  );

  return (
    <View style={styles.container}>
      <View style={styles.topRow}>
        <Badge label={statusLabel} variant={statusVariant} size="md" pulse={game.status === "running"} />
        <Text style={styles.gameCode}>
          {t("results.gameCode", lang)}: {game.code}
        </Text>
      </View>

      <Text style={styles.nickname}>{player.nickname}</Text>

      {team && (
        <View style={styles.teamRow}>
          <View style={[styles.teamDot, { backgroundColor: team.color }]} />
          <Text style={styles.teamName}>
            {t("results.team", lang)}: {team.name}
          </Text>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    alignItems: "center",
    gap: Theme.spacing.sm,
    paddingVertical: Theme.spacing.md,
  },
  topRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: Theme.spacing.md,
  },
  gameCode: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textMuted,
    letterSpacing: 1,
  },
  nickname: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xxxl,
    color: Theme.colors.textPrimary,
    letterSpacing: 2,
    textTransform: "uppercase",
  },
  teamRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: Theme.spacing.sm,
  },
  teamDot: {
    width: 10,
    height: 10,
    borderRadius: 5,
  },
  teamName: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.md,
    color: Theme.colors.textSecondary,
    letterSpacing: 0.5,
  },
});
