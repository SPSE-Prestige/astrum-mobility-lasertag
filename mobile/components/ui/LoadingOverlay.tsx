import React from "react";
import { View, ActivityIndicator, Text, StyleSheet } from "react-native";
import Theme from "@/lib/theme";

interface LoadingOverlayProps {
  message?: string;
}

export default function LoadingOverlay({ message }: LoadingOverlayProps) {
  return (
    <View style={styles.container}>
      <ActivityIndicator size="large" color={Theme.colors.brand} />
      {message && <Text style={styles.text}>{message}</Text>}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: Theme.colors.background,
    gap: Theme.spacing.lg,
  },
  text: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.md,
    color: Theme.colors.textSecondary,
    letterSpacing: 1,
  },
});
