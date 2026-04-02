import React, { useState, useRef, useCallback } from "react";
import {
  View,
  TextInput as RNTextInput,
  StyleSheet,
  Keyboard,
  Platform,
} from "react-native";
import * as Haptics from "expo-haptics";
import Theme from "@/lib/theme";
import { SESSION_CODE_LENGTH } from "@/constants/config";

interface CodeInputProps {
  value: string;
  onChangeText: (text: string) => void;
  onComplete?: (code: string) => void;
  editable?: boolean;
}

export default function CodeInput({
  value,
  onChangeText,
  onComplete,
  editable = true,
}: CodeInputProps) {
  const inputRef = useRef<RNTextInput>(null);

  const handleChange = useCallback(
    (text: string) => {
      const cleaned = text.replace(/[^A-Za-z0-9]/g, "").toUpperCase();
      const limited = cleaned.slice(0, SESSION_CODE_LENGTH);
      onChangeText(limited);

      if (limited.length === SESSION_CODE_LENGTH) {
        Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
        Keyboard.dismiss();
        onComplete?.(limited);
      } else if (limited.length > value.length) {
        Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
      }
    },
    [onChangeText, onComplete, value.length]
  );

  const cells = Array.from({ length: SESSION_CODE_LENGTH }, (_, i) => value[i] ?? "");

  return (
    <View style={styles.container}>
      {/* Visual cells */}
      <View style={styles.cellRow}>
        {cells.map((char, i) => (
          <View
            key={i}
            style={[
              styles.cell,
              char ? styles.cellFilled : null,
              i === value.length && editable ? styles.cellActive : null,
            ]}
          >
            <RNTextInput
              style={styles.cellText}
              value={char}
              editable={false}
            />
          </View>
        ))}
      </View>

      {/* Hidden real input */}
      <RNTextInput
        ref={inputRef}
        style={styles.hiddenInput}
        value={value}
        onChangeText={handleChange}
        maxLength={SESSION_CODE_LENGTH}
        autoCapitalize="characters"
        autoCorrect={false}
        autoFocus
        keyboardType={Platform.OS === "ios" ? "default" : "visible-password"}
        editable={editable}
        caretHidden
        selectionColor="transparent"
        onFocus={() => inputRef.current?.setNativeProps({ selection: { start: value.length, end: value.length } })}
      />

      {/* Tap overlay to focus */}
      <View
        style={StyleSheet.absoluteFill}
        onTouchEnd={() => inputRef.current?.focus()}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    position: "relative",
    alignSelf: "center",
  },
  cellRow: {
    flexDirection: "row",
    gap: Theme.spacing.sm,
  },
  cell: {
    width: 48,
    height: 60,
    borderRadius: Theme.borderRadius.md,
    borderWidth: 1.5,
    borderColor: Theme.colors.border,
    backgroundColor: Theme.colors.surface,
    justifyContent: "center",
    alignItems: "center",
  },
  cellFilled: {
    borderColor: Theme.colors.brand,
    backgroundColor: Theme.colors.surfaceElevated,
  },
  cellActive: {
    borderColor: Theme.colors.brandLight,
    ...Theme.shadows.glow,
  },
  cellText: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xxl,
    color: Theme.colors.textPrimary,
    textAlign: "center",
    letterSpacing: 0,
    padding: 0,
  },
  hiddenInput: {
    position: "absolute",
    width: 1,
    height: 1,
    opacity: 0,
  },
});
