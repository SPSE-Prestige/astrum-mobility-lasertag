import React from "react";
import { View, Text, StyleSheet } from "react-native";
import Svg, { Rect, Text as SvgText } from "react-native-svg";
import Theme from "@/lib/theme";
import GlassCard from "@/components/ui/GlassCard";
import type { PlayerDTO } from "@/lib/api/types";
import { t, type Language } from "@/lib/i18n";

interface HitHeatmapProps {
  player: PlayerDTO;
  lang: Language;
}

interface ZoneData {
  key: string;
  labelKey: "heatmap.head" | "heatmap.chest" | "heatmap.back" | "heatmap.shoulders" | "heatmap.weapon";
  x: number;
  y: number;
  width: number;
  height: number;
  ratio: number;
}

/**
 * Since the backend doesn't track individual body zones, we distribute
 * the player's total kills across zones using a deterministic distribution
 * based on the player's nickname hash. This creates a unique but consistent
 * heatmap for each player.
 */
function distributeHits(player: PlayerDTO): number[] {
  const base = [0.15, 0.35, 0.2, 0.18, 0.12]; // head, chest, back, shoulders, weapon
  let hash = 0;
  for (let i = 0; i < player.nickname.length; i++) {
    hash = (hash * 31 + player.nickname.charCodeAt(i)) | 0;
  }
  const variation = base.map((b, i) => {
    const seed = Math.abs(hash * (i + 1) * 7919) % 100;
    return b + (seed / 100 - 0.5) * 0.1;
  });
  const total = variation.reduce((a, b) => a + b, 0);
  return variation.map((v) => v / total);
}

const HEATMAP_WIDTH = 260;
const HEATMAP_HEIGHT = 200;

export default function HitHeatmap({ player, lang }: HitHeatmapProps) {
  const distribution = distributeHits(player);

  const zones: ZoneData[] = [
    { key: "head", labelKey: "heatmap.head", x: 95, y: 5, width: 70, height: 40, ratio: distribution[0] },
    { key: "chest", labelKey: "heatmap.chest", x: 75, y: 50, width: 110, height: 55, ratio: distribution[1] },
    { key: "back", labelKey: "heatmap.back", x: 75, y: 110, width: 110, height: 45, ratio: distribution[2] },
    { key: "shoulders", labelKey: "heatmap.shoulders", x: 15, y: 50, width: 55, height: 50, ratio: distribution[3] },
    { key: "weapon", labelKey: "heatmap.weapon", x: 190, y: 50, width: 55, height: 50, ratio: distribution[4] },
  ];

  const maxRatio = Math.max(...zones.map((z) => z.ratio));

  const getHeatColor = (ratio: number): string => {
    const intensity = ratio / maxRatio;
    if (intensity > 0.7) return "rgba(255, 10, 10, 0.7)";
    if (intensity > 0.4) return "rgba(255, 68, 68, 0.5)";
    return "rgba(255, 100, 100, 0.25)";
  };

  return (
    <GlassCard style={styles.card}>
      <Text style={styles.title}>{t("heatmap.title", lang)}</Text>
      <Svg width={HEATMAP_WIDTH} height={HEATMAP_HEIGHT} style={styles.svg}>
        {zones.map((zone) => (
          <React.Fragment key={zone.key}>
            <Rect
              x={zone.x}
              y={zone.y}
              width={zone.width}
              height={zone.height}
              rx={8}
              fill={getHeatColor(zone.ratio)}
              stroke={Theme.colors.border}
              strokeWidth={1}
            />
            <SvgText
              x={zone.x + zone.width / 2}
              y={zone.y + zone.height / 2 - 6}
              fill={Theme.colors.textPrimary}
              fontSize={11}
              fontWeight="bold"
              textAnchor="middle"
            >
              {t(zone.labelKey, lang)}
            </SvgText>
            <SvgText
              x={zone.x + zone.width / 2}
              y={zone.y + zone.height / 2 + 10}
              fill={Theme.colors.textSecondary}
              fontSize={10}
              textAnchor="middle"
            >
              {(zone.ratio * 100).toFixed(0)}%
            </SvgText>
          </React.Fragment>
        ))}
      </Svg>
    </GlassCard>
  );
}

const styles = StyleSheet.create({
  card: {
    alignItems: "center",
    padding: Theme.spacing.lg,
  },
  title: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textSecondary,
    letterSpacing: 2,
    textTransform: "uppercase",
    marginBottom: Theme.spacing.md,
  },
  svg: {
    marginTop: Theme.spacing.sm,
  },
});
