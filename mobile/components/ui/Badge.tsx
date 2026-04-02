import React from "react";
import { View, Text, StyleSheet, type ViewStyle } from "react-native";
import Theme from "@/lib/theme";

type BadgeVariant = "brand" | "success" | "warning" | "danger" | "info" | "muted";

interface BadgeProps {
  label: string;
  variant?: BadgeVariant;
  size?: "sm" | "md";
  pulse?: boolean;
  style?: ViewStyle;
}

const variantColors: Record<BadgeVariant, { bg: string; text: string; dot: string }> = {
  brand: {
    bg: Theme.colors.brandMuted,
    text: Theme.colors.brand,
    dot: Theme.colors.brand,
  },
  success: {
    bg: Theme.colors.successMuted,
    text: Theme.colors.success,
    dot: Theme.colors.success,
  },
  warning: {
    bg: Theme.colors.warningMuted,
    text: Theme.colors.warning,
    dot: Theme.colors.warning,
  },
  danger: {
    bg: Theme.colors.dangerMuted,
    text: Theme.colors.danger,
    dot: Theme.colors.danger,
  },
  info: {
    bg: Theme.colors.infoMuted,
    text: Theme.colors.info,
    dot: Theme.colors.info,
  },
  muted: {
    bg: Theme.colors.surfaceElevated,
    text: Theme.colors.textSecondary,
    dot: Theme.colors.textMuted,
  },
};

export default function Badge({
  label,
  variant = "brand",
  size = "sm",
  pulse = false,
  style,
}: BadgeProps) {
  const colors = variantColors[variant];

  return (
    <View
      style={[
        styles.badge,
        size === "md" && styles.badgeMd,
        { backgroundColor: colors.bg },
        style,
      ]}
    >
      {pulse && (
        <View style={[styles.dot, { backgroundColor: colors.dot }]} />
      )}
      <Text
        style={[
          styles.text,
          size === "md" && styles.textMd,
          { color: colors.text },
        ]}
      >
        {label}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: Theme.spacing.sm,
    paddingVertical: Theme.spacing.xxs,
    borderRadius: Theme.borderRadius.full,
    gap: Theme.spacing.xs,
  },
  badgeMd: {
    paddingHorizontal: Theme.spacing.md,
    paddingVertical: Theme.spacing.xs,
  },
  dot: {
    width: 6,
    height: 6,
    borderRadius: 3,
  },
  text: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xs,
    letterSpacing: 1,
    textTransform: "uppercase",
  },
  textMd: {
    fontSize: Theme.fontSize.sm,
  },
});
