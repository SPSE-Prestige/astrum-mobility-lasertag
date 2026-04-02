import React from "react";
import {
  View,
  TextInput as RNTextInput,
  Text,
  StyleSheet,
  type ViewStyle,
} from "react-native";
import Theme from "@/lib/theme";

interface TextInputProps {
  value: string;
  onChangeText: (text: string) => void;
  placeholder?: string;
  label?: string;
  error?: string | null;
  autoCapitalize?: "none" | "sentences" | "words" | "characters";
  autoFocus?: boolean;
  maxLength?: number;
  keyboardType?: "default" | "number-pad" | "email-address";
  style?: ViewStyle;
  onSubmitEditing?: () => void;
}

export default function TextInput({
  value,
  onChangeText,
  placeholder,
  label,
  error,
  autoCapitalize = "characters",
  autoFocus = false,
  maxLength,
  keyboardType = "default",
  style,
  onSubmitEditing,
}: TextInputProps) {
  return (
    <View style={[styles.container, style]}>
      {label && <Text style={styles.label}>{label}</Text>}
      <RNTextInput
        style={[styles.input, error ? styles.inputError : null]}
        value={value}
        onChangeText={onChangeText}
        placeholder={placeholder}
        placeholderTextColor={Theme.colors.textMuted}
        autoCapitalize={autoCapitalize}
        autoFocus={autoFocus}
        maxLength={maxLength}
        keyboardType={keyboardType}
        selectionColor={Theme.colors.brand}
        returnKeyType="done"
        onSubmitEditing={onSubmitEditing}
      />
      {error && <Text style={styles.error}>{error}</Text>}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    width: "100%",
  },
  label: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.textSecondary,
    marginBottom: Theme.spacing.sm,
    letterSpacing: 1,
    textTransform: "uppercase",
  },
  input: {
    backgroundColor: Theme.colors.surface,
    borderWidth: 1.5,
    borderColor: Theme.colors.border,
    borderRadius: Theme.borderRadius.lg,
    paddingHorizontal: Theme.spacing.lg,
    paddingVertical: Theme.spacing.md,
    fontSize: Theme.fontSize.xl,
    fontFamily: Theme.fontFamily.display,
    color: Theme.colors.textPrimary,
    textAlign: "center",
    letterSpacing: 8,
    height: 60,
  },
  inputError: {
    borderColor: Theme.colors.danger,
  },
  error: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.danger,
    marginTop: Theme.spacing.sm,
    textAlign: "center",
  },
});
