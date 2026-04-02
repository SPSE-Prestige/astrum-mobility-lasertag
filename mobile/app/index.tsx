import React, { useCallback, useContext, useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  TouchableOpacity,
  ScrollView,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { router } from "expo-router";
import Theme from "@/lib/theme";
import { t } from "@/lib/i18n";
import { LanguageContext } from "./_layout";
import CodeInput from "@/components/login/CodeInput";
import Button from "@/components/ui/Button";
import usePlayerSession from "@/hooks/usePlayerSession";
import { SESSION_CODE_LENGTH } from "@/constants/config";

export default function LoginScreen() {
  const { language, setLanguage } = useContext(LanguageContext);
  const [code, setCode] = useState("");
  const session = usePlayerSession();

  const handleSubmit = useCallback(
    async (submittedCode?: string) => {
      const finalCode = submittedCode ?? code;
      if (finalCode.length !== SESSION_CODE_LENGTH) return;

      await session.fetch(finalCode);
    },
    [code, session]
  );

  // Navigate on successful fetch
  React.useEffect(() => {
    if (session.data) {
      router.push({
        pathname: "/results",
        params: { code: code.toUpperCase() },
      });
    }
  }, [session.data, code]);

  const errorMessage = session.error
    ? session.errorCode === "NOT_FOUND"
      ? t("login.error.notFound", language)
      : session.errorCode === "NETWORK_ERROR"
        ? t("login.error.network", language)
        : t("login.error.generic", language)
    : null;

  return (
    <SafeAreaView style={styles.safe} edges={["top", "bottom"]}>
      <KeyboardAvoidingView
        style={styles.flex}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        <ScrollView
          contentContainerStyle={styles.scroll}
          keyboardShouldPersistTaps="handled"
        >
          {/* Language Switcher */}
          <View style={styles.langRow}>
            {(["cs", "en"] as const).map((lang) => (
              <TouchableOpacity
                key={lang}
                onPress={() => setLanguage(lang)}
                style={[styles.langBtn, language === lang && styles.langBtnActive]}
              >
                <Text
                  style={[
                    styles.langText,
                    language === lang && styles.langTextActive,
                  ]}
                >
                  {t(`lang.${lang}` as "lang.cs" | "lang.en", language)}
                </Text>
              </TouchableOpacity>
            ))}
          </View>

          {/* Brand */}
          <View style={styles.brand}>
            <View style={styles.logoContainer}>
              <Text style={styles.logoIcon}>⎊</Text>
            </View>
            <Text style={styles.title}>{t("login.title", language)}</Text>
            <Text style={styles.subtitle}>{t("login.subtitle", language)}</Text>
          </View>

          {/* Code Input */}
          <View style={styles.inputSection}>
            <CodeInput
              value={code}
              onChangeText={(text) => {
                setCode(text);
                session.clear();
              }}
              onComplete={handleSubmit}
              editable={!session.loading}
            />

            {errorMessage && (
              <Text style={styles.error}>{errorMessage}</Text>
            )}

            <Button
              title={
                session.loading
                  ? t("login.loading", language)
                  : t("login.button", language)
              }
              onPress={() => handleSubmit()}
              loading={session.loading}
              disabled={code.length !== SESSION_CODE_LENGTH}
              style={styles.button}
              size="lg"
            />
          </View>

          {/* Footer */}
          <View style={styles.footer}>
            <Text style={styles.footerText}>ASTRUM MOBILITY</Text>
            <View style={styles.footerDot} />
            <Text style={styles.footerText}>LASER TAG</Text>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: {
    flex: 1,
    backgroundColor: Theme.colors.background,
  },
  flex: {
    flex: 1,
  },
  scroll: {
    flexGrow: 1,
    paddingHorizontal: Theme.spacing.xxl,
    justifyContent: "center",
  },
  langRow: {
    flexDirection: "row",
    justifyContent: "center",
    gap: Theme.spacing.sm,
    position: "absolute",
    top: Theme.spacing.lg,
    left: 0,
    right: 0,
    zIndex: 10,
  },
  langBtn: {
    paddingHorizontal: Theme.spacing.md,
    paddingVertical: Theme.spacing.xs,
    borderRadius: Theme.borderRadius.sm,
    borderWidth: 1,
    borderColor: Theme.colors.border,
  },
  langBtnActive: {
    borderColor: Theme.colors.brand,
    backgroundColor: Theme.colors.brandMuted,
  },
  langText: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xs,
    color: Theme.colors.textMuted,
    letterSpacing: 1,
  },
  langTextActive: {
    color: Theme.colors.brand,
  },
  brand: {
    alignItems: "center",
    marginBottom: Theme.spacing.huge,
  },
  logoContainer: {
    width: 72,
    height: 72,
    borderRadius: 36,
    backgroundColor: Theme.colors.brandMuted,
    justifyContent: "center",
    alignItems: "center",
    marginBottom: Theme.spacing.xl,
    borderWidth: 2,
    borderColor: Theme.colors.brand,
    ...Theme.shadows.glow,
  },
  logoIcon: {
    fontSize: 32,
    color: Theme.colors.brand,
  },
  title: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xxxl,
    color: Theme.colors.textPrimary,
    letterSpacing: 4,
    textAlign: "center",
  },
  subtitle: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.md,
    color: Theme.colors.textMuted,
    marginTop: Theme.spacing.sm,
    textAlign: "center",
    lineHeight: 22,
  },
  inputSection: {
    alignItems: "center",
    gap: Theme.spacing.xl,
  },
  error: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.sm,
    color: Theme.colors.danger,
    textAlign: "center",
  },
  button: {
    width: "100%",
    maxWidth: 320,
  },
  footer: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    marginTop: Theme.spacing.huge,
    gap: Theme.spacing.sm,
  },
  footerText: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xs,
    color: Theme.colors.textMuted,
    letterSpacing: 2,
  },
  footerDot: {
    width: 4,
    height: 4,
    borderRadius: 2,
    backgroundColor: Theme.colors.brand,
  },
});
