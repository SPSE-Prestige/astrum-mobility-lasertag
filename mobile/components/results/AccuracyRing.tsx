import React from "react";
import { View, Text, StyleSheet } from "react-native";
import Svg, { Circle } from "react-native-svg";
import Theme from "@/lib/theme";
import GlassCard from "@/components/ui/GlassCard";
import { t, type Language } from "@/lib/i18n";

interface AccuracyRingProps {
  kills: number;
  shotsFired: number;
  lang: Language;
}

const RING_SIZE = 140;
const STROKE_WIDTH = 10;
const RADIUS = (RING_SIZE - STROKE_WIDTH) / 2;
const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

export default function AccuracyRing({ kills, shotsFired, lang }: AccuracyRingProps) {
  const accuracy = shotsFired > 0 ? (kills / shotsFired) * 100 : 0;
  const progress = Math.min(accuracy / 100, 1);
  const strokeDashoffset = CIRCUMFERENCE * (1 - progress);

  const getColor = () => {
    if (accuracy >= 50) return Theme.colors.success;
    if (accuracy >= 25) return Theme.colors.warning;
    return Theme.colors.danger;
  };

  return (
    <GlassCard style={styles.card}>
      <Text style={styles.title}>{t("stats.accuracy", lang)}</Text>
      <View style={styles.ringContainer}>
        <Svg width={RING_SIZE} height={RING_SIZE}>
          {/* Background ring */}
          <Circle
            cx={RING_SIZE / 2}
            cy={RING_SIZE / 2}
            r={RADIUS}
            stroke={Theme.colors.border}
            strokeWidth={STROKE_WIDTH}
            fill="none"
          />
          {/* Progress ring */}
          <Circle
            cx={RING_SIZE / 2}
            cy={RING_SIZE / 2}
            r={RADIUS}
            stroke={getColor()}
            strokeWidth={STROKE_WIDTH}
            fill="none"
            strokeLinecap="round"
            strokeDasharray={CIRCUMFERENCE}
            strokeDashoffset={strokeDashoffset}
            rotation="-90"
            origin={`${RING_SIZE / 2}, ${RING_SIZE / 2}`}
          />
        </Svg>
        <View style={styles.centerLabel}>
          <Text style={[styles.percentage, { color: getColor() }]}>
            {accuracy.toFixed(1)}
          </Text>
          <Text style={styles.percentSign}>%</Text>
        </View>
      </View>
      <View style={styles.details}>
        <Text style={styles.detail}>
          {kills} {t("stats.kills", lang).toLowerCase()} / {shotsFired} {t("stats.shots", lang).toLowerCase()}
        </Text>
      </View>
    </GlassCard>
  );
}

const styles = StyleSheet.create({
  card: {
    alignItems: "center",
    padding: Theme.spacing.xl,
  },
  title: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textSecondary,
    letterSpacing: 2,
    textTransform: "uppercase",
    marginBottom: Theme.spacing.lg,
  },
  ringContainer: {
    width: RING_SIZE,
    height: RING_SIZE,
    justifyContent: "center",
    alignItems: "center",
  },
  centerLabel: {
    position: "absolute",
    alignItems: "center",
  },
  percentage: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xxxl,
    lineHeight: Theme.fontSize.xxxl * 1.1,
  },
  percentSign: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textMuted,
    marginTop: -2,
  },
  details: {
    marginTop: Theme.spacing.md,
  },
  detail: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textMuted,
    textAlign: "center",
  },
});
