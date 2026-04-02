import React, { useCallback, useEffect, useState } from "react";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useFonts } from "expo-font";
import * as SplashScreen from "expo-splash-screen";
import { View, StyleSheet } from "react-native";
import Theme from "@/lib/theme";
import type { Language } from "@/lib/i18n";

SplashScreen.preventAutoHideAsync();

export const LanguageContext = React.createContext<{
  language: Language;
  setLanguage: (lang: Language) => void;
}>({
  language: "cs",
  setLanguage: () => {},
});

export default function RootLayout() {
  const [language, setLanguage] = useState<Language>("cs");

  const [fontsLoaded] = useFonts({
    "Goldman-Regular": require("@/assets/fonts/Goldman-Regular.ttf"),
    "Goldman-Bold": require("@/assets/fonts/Goldman-Bold.ttf"),
  });

  const onLayoutRootView = useCallback(async () => {
    if (fontsLoaded) {
      await SplashScreen.hideAsync();
    }
  }, [fontsLoaded]);

  if (!fontsLoaded) return null;

  return (
    <LanguageContext.Provider value={{ language, setLanguage }}>
      <View style={styles.root} onLayout={onLayoutRootView}>
        <StatusBar style="light" />
        <Stack
          screenOptions={{
            headerShown: false,
            contentStyle: { backgroundColor: Theme.colors.background },
            animation: "slide_from_right",
          }}
        />
      </View>
    </LanguageContext.Provider>
  );
}

const styles = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: Theme.colors.background,
  },
});
