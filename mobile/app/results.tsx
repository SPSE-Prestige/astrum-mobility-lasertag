import React, { useCallback, useContext, useEffect } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  RefreshControl,
  TouchableOpacity,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { router, useLocalSearchParams } from "expo-router";
import Theme from "@/lib/theme";
import { t } from "@/lib/i18n";
import { LanguageContext } from "./_layout";
import usePlayerSession from "@/hooks/usePlayerSession";
import usePolling from "@/hooks/usePolling";
import LoadingOverlay from "@/components/ui/LoadingOverlay";
import PlayerHeader from "@/components/results/PlayerHeader";
import HeroScore from "@/components/results/HeroScore";
import StatsGrid from "@/components/results/StatsGrid";
import AccuracyRing from "@/components/results/AccuracyRing";
import HitHeatmap from "@/components/results/HitHeatmap";
import Leaderboard from "@/components/results/Leaderboard";
import KillFeed from "@/components/results/KillFeed";

export default function ResultsScreen() {
  const { language, setLanguage } = useContext(LanguageContext);
  const params = useLocalSearchParams<{ code: string }>();
  const session = usePlayerSession();
  const code = params.code ?? "";

  // Initial fetch
  useEffect(() => {
    if (code) {
      session.fetch(code);
    }
  }, [code]);

  // Auto-poll for live games
  usePolling({
    enabled: session.data?.game.status === "running",
    onTick: session.refresh,
  });

  const handleBack = useCallback(() => {
    session.clear();
    router.back();
  }, [session]);

  if (session.loading && !session.data) {
    return <LoadingOverlay message={t("common.loading", language)} />;
  }

  if (!session.data) {
    return (
      <SafeAreaView style={styles.safe}>
        <View style={styles.errorContainer}>
          <Text style={styles.errorText}>
            {session.error ?? t("common.noData", language)}
          </Text>
          <TouchableOpacity onPress={handleBack} style={styles.backButton}>
            <Text style={styles.backButtonText}>{t("common.back", language)}</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  const { player, game, team, leaderboard, events } = session.data;

  return (
    <SafeAreaView style={styles.safe} edges={["top", "bottom"]}>
      <ScrollView
        style={styles.scroll}
        contentContainerStyle={styles.content}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={session.loading}
            onRefresh={session.refresh}
            tintColor={Theme.colors.brand}
            colors={[Theme.colors.brand]}
          />
        }
      >
        {/* Top Bar */}
        <View style={styles.topBar}>
          <TouchableOpacity onPress={handleBack} style={styles.backBtn}>
            <Text style={styles.backBtnText}>← {t("common.back", language)}</Text>
          </TouchableOpacity>
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
        </View>

        {/* Player Header */}
        <PlayerHeader player={player} game={game} team={team} lang={language} />

        {/* Hero Score */}
        <HeroScore score={player.score} lang={language} />

        {/* Stats Grid */}
        <View style={styles.section}>
          <StatsGrid player={player} lang={language} />
        </View>

        {/* Accuracy Ring + Heatmap side by side on wide screens, stacked on narrow */}
        <View style={styles.dualRow}>
          <View style={styles.dualItem}>
            <AccuracyRing
              kills={player.kills}
              shotsFired={player.shots_fired}
              lang={language}
            />
          </View>
          <View style={styles.dualItem}>
            <HitHeatmap player={player} lang={language} />
          </View>
        </View>

        {/* Weapon Level */}
        <View style={styles.weaponRow}>
          <View style={styles.weaponCard}>
            <Text style={styles.weaponLabel}>{t("stats.weaponLevel", language)}</Text>
            <View style={styles.weaponLevelRow}>
              {Array.from({ length: Math.max(player.weapon_level, 1) }, (_, i) => (
                <View
                  key={i}
                  style={[
                    styles.weaponDot,
                    i < player.weapon_level && styles.weaponDotActive,
                  ]}
                />
              ))}
              <Text style={styles.weaponValue}>LVL {player.weapon_level}</Text>
            </View>
          </View>
        </View>

        {/* Leaderboard */}
        {leaderboard && leaderboard.length > 0 && (
          <View style={styles.section}>
            <Leaderboard players={leaderboard} lang={language} />
          </View>
        )}

        {/* Kill Feed */}
        {events && events.length > 0 && (
          <View style={styles.section}>
            <KillFeed events={events} currentNickname={player.nickname} lang={language} />
          </View>
        )}

        {/* Bottom spacing */}
        <View style={styles.bottomSpacer} />
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: {
    flex: 1,
    backgroundColor: Theme.colors.background,
  },
  scroll: {
    flex: 1,
  },
  content: {
    paddingHorizontal: Theme.spacing.lg,
    gap: Theme.spacing.lg,
  },
  topBar: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingTop: Theme.spacing.sm,
  },
  backBtn: {
    paddingVertical: Theme.spacing.sm,
    paddingRight: Theme.spacing.md,
  },
  backBtnText: {
    fontFamily: Theme.fontFamily.displayRegular,
    fontSize: Theme.fontSize.md,
    color: Theme.colors.textSecondary,
  },
  langRow: {
    flexDirection: "row",
    gap: Theme.spacing.xs,
  },
  langBtn: {
    paddingHorizontal: Theme.spacing.sm,
    paddingVertical: Theme.spacing.xxs,
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
  section: {
    gap: Theme.spacing.sm,
  },
  dualRow: {
    flexDirection: "row",
    gap: Theme.spacing.sm,
  },
  dualItem: {
    flex: 1,
  },
  weaponRow: {
    alignItems: "center",
  },
  weaponCard: {
    backgroundColor: Theme.colors.surface,
    borderRadius: Theme.borderRadius.lg,
    borderWidth: 1,
    borderColor: Theme.colors.border,
    paddingHorizontal: Theme.spacing.xl,
    paddingVertical: Theme.spacing.md,
    alignItems: "center",
    gap: Theme.spacing.sm,
  },
  weaponLabel: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.xs,
    color: Theme.colors.textMuted,
    letterSpacing: 2,
    textTransform: "uppercase",
  },
  weaponLevelRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: Theme.spacing.xs,
  },
  weaponDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: Theme.colors.border,
  },
  weaponDotActive: {
    backgroundColor: Theme.colors.warning,
  },
  weaponValue: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.lg,
    color: Theme.colors.warning,
    marginLeft: Theme.spacing.sm,
  },
  errorContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: Theme.spacing.xxl,
    gap: Theme.spacing.xl,
  },
  errorText: {
    fontFamily: Theme.fontFamily.body,
    fontSize: Theme.fontSize.md,
    color: Theme.colors.danger,
    textAlign: "center",
  },
  backButton: {
    paddingHorizontal: Theme.spacing.xl,
    paddingVertical: Theme.spacing.md,
    borderRadius: Theme.borderRadius.lg,
    borderWidth: 1,
    borderColor: Theme.colors.brand,
  },
  backButtonText: {
    fontFamily: Theme.fontFamily.display,
    fontSize: Theme.fontSize.md,
    color: Theme.colors.brand,
    letterSpacing: 1,
  },
  bottomSpacer: {
    height: Theme.spacing.xxl,
  },
});
