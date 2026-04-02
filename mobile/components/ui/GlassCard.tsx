import React, { type ReactNode } from "react";
import { View, StyleSheet, type ViewStyle } from "react-native";
import Theme from "@/lib/theme";

interface GlassCardProps {
  children: ReactNode;
  style?: ViewStyle;
  glow?: boolean;
  padding?: keyof typeof Theme.spacing;
}

export default function GlassCard({
  children,
  style,
  glow = false,
  padding = "lg",
}: GlassCardProps) {
  return (
    <View
      style={[
        styles.card,
        { padding: Theme.spacing[padding] },
        glow && styles.glow,
        style,
      ]}
    >
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: "rgba(17, 17, 17, 0.85)",
    borderRadius: Theme.borderRadius.xl,
    borderWidth: 1,
    borderColor: Theme.colors.border,
    ...Theme.shadows.md,
  },
  glow: {
    borderColor: Theme.colors.brandMuted,
    ...Theme.shadows.glow,
  },
});
