import { Platform } from "react-native";

const Colors = {
  background: "#0a0a0a",
  surface: "#111111",
  surfaceElevated: "#1a1a1a",
  surfaceBright: "#222222",
  border: "#2a2a2a",
  borderLight: "#333333",

  textPrimary: "#ffffff",
  textSecondary: "#a0a0a0",
  textMuted: "#666666",
  textInverse: "#0a0a0a",

  brand: "#ff0a0a",
  brandLight: "#ff4444",
  brandDark: "#cc0000",
  brandMuted: "#ff0a0a33",
  brandGlow: "rgba(255, 10, 10, 0.12)",

  success: "#22c55e",
  successMuted: "#22c55e22",
  warning: "#f59e0b",
  warningMuted: "#f59e0b22",
  info: "#3b82f6",
  infoMuted: "#3b82f622",
  danger: "#ef4444",
  dangerMuted: "#ef444422",

  white: "#ffffff",
  black: "#000000",
  transparent: "transparent",
} as const;

const Spacing = {
  xxs: 2,
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 20,
  xxl: 24,
  xxxl: 32,
  huge: 48,
  massive: 64,
} as const;

const BorderRadius = {
  sm: 6,
  md: 10,
  lg: 14,
  xl: 18,
  xxl: 24,
  full: 9999,
} as const;

const FontFamily = {
  display: "Goldman-Bold",
  displayRegular: "Goldman-Regular",
  body: Platform.select({
    ios: "System",
    android: "Roboto",
    default: "System",
  }) as string,
} as const;

const FontSize = {
  xs: 11,
  sm: 13,
  md: 15,
  lg: 17,
  xl: 20,
  xxl: 24,
  xxxl: 32,
  hero: 48,
  mega: 64,
} as const;

const Shadows = {
  sm: {
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.2,
    shadowRadius: 2,
    elevation: 2,
  },
  md: {
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 4,
  },
  lg: {
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.4,
    shadowRadius: 16,
    elevation: 8,
  },
  glow: {
    shadowColor: Colors.brand,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.5,
    shadowRadius: 20,
    elevation: 6,
  },
} as const;

export const Theme = {
  colors: Colors,
  spacing: Spacing,
  borderRadius: BorderRadius,
  fontFamily: FontFamily,
  fontSize: FontSize,
  shadows: Shadows,
} as const;

export type ThemeType = typeof Theme;
export default Theme;
