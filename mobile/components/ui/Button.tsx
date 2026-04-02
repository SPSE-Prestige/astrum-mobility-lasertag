import React from "react";
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  ActivityIndicator,
  type ViewStyle,
  type TextStyle,
} from "react-native";
import * as Haptics from "expo-haptics";
import Theme from "@/lib/theme";

interface ButtonProps {
  title: string;
  onPress: () => void;
  variant?: "primary" | "secondary" | "ghost";
  size?: "sm" | "md" | "lg";
  loading?: boolean;
  disabled?: boolean;
  style?: ViewStyle;
  textStyle?: TextStyle;
}

export default function Button({
  title,
  onPress,
  variant = "primary",
  size = "md",
  loading = false,
  disabled = false,
  style,
  textStyle,
}: ButtonProps) {
  const handlePress = () => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);
    onPress();
  };

  const isDisabled = disabled || loading;

  return (
    <TouchableOpacity
      style={[
        styles.base,
        styles[variant],
        styles[`size_${size}`],
        isDisabled && styles.disabled,
        style,
      ]}
      onPress={handlePress}
      disabled={isDisabled}
      activeOpacity={0.7}
    >
      {loading ? (
        <ActivityIndicator
          size="small"
          color={variant === "primary" ? Theme.colors.white : Theme.colors.brand}
        />
      ) : (
        <Text
          style={[
            styles.text,
            styles[`text_${variant}`],
            styles[`textSize_${size}`],
            textStyle,
          ]}
        >
          {title}
        </Text>
      )}
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  base: {
    alignItems: "center",
    justifyContent: "center",
    borderRadius: Theme.borderRadius.lg,
    flexDirection: "row",
  },
  primary: {
    backgroundColor: Theme.colors.brand,
    ...Theme.shadows.glow,
  },
  secondary: {
    backgroundColor: Theme.colors.transparent,
    borderWidth: 1.5,
    borderColor: Theme.colors.brand,
  },
  ghost: {
    backgroundColor: Theme.colors.transparent,
  },
  size_sm: {
    height: 36,
    paddingHorizontal: Theme.spacing.lg,
  },
  size_md: {
    height: 48,
    paddingHorizontal: Theme.spacing.xl,
  },
  size_lg: {
    height: 56,
    paddingHorizontal: Theme.spacing.xxl,
  },
  disabled: {
    opacity: 0.5,
  },
  text: {
    fontFamily: Theme.fontFamily.display,
    letterSpacing: 1.5,
  },
  text_primary: {
    color: Theme.colors.white,
  },
  text_secondary: {
    color: Theme.colors.brand,
  },
  text_ghost: {
    color: Theme.colors.textSecondary,
  },
  textSize_sm: {
    fontSize: Theme.fontSize.sm,
  },
  textSize_md: {
    fontSize: Theme.fontSize.md,
  },
  textSize_lg: {
    fontSize: Theme.fontSize.lg,
  },
});
