import React from "react";
import { View, Text, StyleSheet } from "react-native";
import Theme from "@/lib/theme";
import GlassCard from "@/components/ui/GlassCard";
import { t, type Language } from "@/lib/i18n";

interface HeroScoreProps {
  score: number;
  lang: Language;
}

export default function HeroScore({ score, lang }: HeroScoreProps) {
  return (
    <GlassCard glow style={styles.card}>
      <Text style={styles.label}>{t("stats.score", lang)}</Text>
      <Text style={styles.score}>{score.toLocaleString()}</Text>
      <View style={styles.underline} />
    </GlassCard>
  );
}

const styles = StyleSheet.create({
  card: {
    alignItems: "center",
    paddingVertical: Theme.spacing.xxl,
  },
  label: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.brand,
    letterSpacing: 3,
    textTransform: "uppercase",
    marginBottom: Theme.spacing.xs,
  },
  score: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.mega,
    color: Theme.colors.textPrimary,
    letterSpacing: 2,
    lineHeight: Theme.fontSize.mega * 1.1,
  },
  underline: {
    width: 60,
    height: 3,
    backgroundColor: Theme.colors.brand,
    borderRadius: 2,
    marginTop: Theme.spacing.md,
  },
});
