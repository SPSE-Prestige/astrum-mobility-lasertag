import { StatusBar } from 'expo-status-bar';
import React, { useState } from 'react';
import { StyleSheet, Text, View, ScrollView, TouchableOpacity, SafeAreaView, Platform, TextInput, Alert, KeyboardAvoidingView } from 'react-native';
import { useFonts } from 'expo-font';
import { Goldman_400Regular, Goldman_700Bold } from '@expo-google-fonts/goldman';


// --- MOCK DATA (Updated based on requirements) ---
const GAME_INFO = {
  mode: 'TEAM DEATHMATCH',
  map: 'ARENA 1',
  duration: '10:00',
  date: '21.03.2026',
};

// Mock database mapping codes to users
const MOCK_USERS_DB: Record<string, { id: string, name: string }> = {
  '1234': { id: 'u1', name: 'Player One' },
  '0007': { id: 'u2', name: 'Terminator' },
  '1111': { id: 'u3', name: 'Sniper_CZ' },
};

const MOCK_CURRENT_USER_TEMPLATE = {
  id: 'u1',
                  // Name will be filled from login
  rank: 2,
  kills: 12,        // Kolikrát sestřelil
  deaths: 8,        // Kolikrát byl sestřelen
  hits: 45,         // Trefené cíle (pro výpočet přesnosti)
  shotsFired: 120,  // Vystřelené náboje (pro výpočet přesnosti)
  favoriteWeapon: 'Plasma Rifle', // Oblíbená zbraň
};

const MOCK_LEADERBOARD = [
  { id: 'u2', rank: 1, name: 'Terminator', kills: 18, deaths: 2, hits: 55 },
  { id: 'u1', rank: 2, name: 'Player One', kills: 12, deaths: 8, hits: 45 },
  { id: 'u3', rank: 3, name: 'Sniper_CZ', kills: 8, deaths: 5, hits: 32 },
  { id: 'u4', rank: 4, name: 'Speedy', kills: 5, deaths: 12, hits: 25 },
  { id: 'u5', rank: 5, name: 'Tank', kills: 2, deaths: 15, hits: 18 },
];

const MOCK_HISTORY = [
  { id: 'h1', date: '20.03.2026', mode: 'FFA', rank: 1, kills: 15, deaths: 5, map: 'ARENA 2' },
  { id: 'h2', date: '18.03.2026', mode: 'TDM', rank: 3, kills: 8, deaths: 10, map: 'FOREST' },
  { id: 'h3', date: '15.03.2026', mode: 'CAPTURE', rank: 2, kills: 10, deaths: 4, map: 'FACTORY' },
];

/**
 * Calculates accuracy percentage based on hits and shots fired.
 */
const calculateAccuracy = (hits: number, shots: number) => {
  if (shots === 0) return 0;
  return Math.round((hits / shots) * 100);
};

// --- THEME CONSTANTS (From web/app/globals.css) ---
const COLORS = {
  background: '#020303', // Zinc-950 equivalent
  cardBackground: '#18181b', // Zinc-900 equivalent
  textPrimary: '#ebedef', // Zinc-200
  textSecondary: '#a1a1aa', // Zinc-400
  accentGreen: '#00ff00',
  accentRed: '#ff0000',
  border: '#27272a', // Zinc-800
};

export default function App() {
  const [fontsLoaded] = useFonts({
    Goldman_400Regular,
    Goldman_700Bold,
  });

  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loginCode, setLoginCode] = useState('');
  const [currentUser, setCurrentUser] = useState(MOCK_CURRENT_USER_TEMPLATE);

  const handleLogin = () => {
    if (MOCK_USERS_DB[loginCode]) {
      setCurrentUser({
        ...MOCK_CURRENT_USER_TEMPLATE,
        name: MOCK_USERS_DB[loginCode].name,
        id: MOCK_USERS_DB[loginCode].id
      });
      setIsLoggedIn(true);
    } else {
      Alert.alert('Chyba', 'Neplatný přihlašovací kód');
    }
  };

  const handleLogout = () => {
    setIsLoggedIn(false);
    setLoginCode('');
  };

  const accuracy = calculateAccuracy(currentUser.hits, currentUser.shotsFired);

  if (!fontsLoaded) {
    return <View style={styles.loadingContainer}><Text style={{color: 'white'}}>Loading...</Text></View>;
  }

  // --- LOGIN SCREEN ---
  if (!isLoggedIn) {
    return (
      <SafeAreaView style={styles.container}>
        <StatusBar style="light" backgroundColor={COLORS.background} />
        <KeyboardAvoidingView 
          behavior={Platform.OS === "ios" ? "padding" : "height"}
          style={styles.loginContainer}
        >
          <Text style={styles.loginLogo}>LASER TAG</Text>
          <Text style={styles.loginSubtitle}>PŘIHLÁŠENÍ HRÁČE</Text>
          
          <View style={styles.inputContainer}>
            <Text style={styles.inputLabel}>ZADEJTE VÁŠ PŘIDĚLENÝ KÓD</Text>
            <TextInput
              style={styles.input}
              placeholder="0000"
              placeholderTextColor={COLORS.textSecondary}
              keyboardType="numeric"
              maxLength={4}
              value={loginCode}
              onChangeText={setLoginCode}
              secureTextEntry
            />
          </View>

          <TouchableOpacity style={styles.loginButton} onPress={handleLogin}>
            <Text style={styles.actionButtonText}>VSTOUPIT DO ARÉNY</Text>
          </TouchableOpacity>
          
          <Text style={styles.loginFooter}>Kód vám přidělí administrátor při registraci</Text>
        </KeyboardAvoidingView>
      </SafeAreaView>
    );
  }

  // --- GAME STATS SCREEN ---
  return (
    <SafeAreaView style={styles.container}>
      <StatusBar style="light" backgroundColor={COLORS.background} />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        
        <View style={styles.topBar}>
           <TouchableOpacity onPress={handleLogout} style={styles.logoutButton}>
             <Text style={styles.logoutText}>ODHLÁSIT</Text>
           </TouchableOpacity>
        </View>

        {/* HEADER - GAME INFO */}
        <View style={styles.header}>
          <Text style={styles.headerTitle}>VÝSLEDKY HRY</Text>
          <View style={styles.gameInfoContainer}>
             <Text style={styles.gameInfoText}>{GAME_INFO.mode}</Text>
             <Text style={styles.separator}>•</Text>
             <Text style={styles.gameInfoText}>{GAME_INFO.map}</Text>
          </View>
          <Text style={styles.dateText}>{GAME_INFO.date} • {GAME_INFO.duration}</Text>
        </View>

        {/* PLAYER CARD (Nickname) */}
         <View style={styles.playerCard}>
          <Text style={styles.playerLabel}>HRÁČ</Text>
          <Text style={styles.playerName}>{currentUser.name}</Text>
          <View style={styles.rankBadge}>
            <Text style={styles.rankText}>#{currentUser.rank} MÍSTO</Text>
          </View>
        </View>

        {/* STATS GRID */}
        <View style={styles.gridContainer}>
          {/* Kills - "Kolik sestřelil" */}
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>SESTŘEILY</Text>
            <Text style={[styles.gridValue, { color: COLORS.accentGreen }]}>{currentUser.kills}</Text>
            <Text style={styles.gridSubtext}>Zadané zásahy</Text>
          </View>

          {/* Deaths - "Kolikrát byl sestřelen" */}
           <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>BYL SESTŘELEN</Text>
            <Text style={[styles.gridValue, { color: COLORS.accentRed }]}>{currentUser.deaths}x</Text>
            <Text style={styles.gridSubtext}>Smrti</Text>
          </View>
          
          {/* Accuracy - "Přesnost" */}
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>PŘESNOST</Text>
            <Text style={styles.gridValue}>{accuracy}%</Text>
            <Text style={styles.gridSubtext}>{currentUser.hits} tref / {currentUser.shotsFired} střel</Text>
          </View>

          {/* Favorite Weapon - "Oblíbená zbraň" */}
          <View style={styles.gridItem}>
            <Text style={styles.gridLabel}>ZBRAŇ</Text>
            <Text style={[styles.gridValue, { fontSize: 20 }]} numberOfLines={1}>{currentUser.favoriteWeapon}</Text>
            <Text style={styles.gridSubtext}>Nejpoužívanější</Text>
          </View>
        </View>

        {/* LEADERBOARD SECTION */}
        <View style={styles.sectionContainer}>
          <Text style={styles.sectionTitle}>ŽEBŘÍČEK</Text>
          <View style={styles.leaderboardContainer}>
            {/* Header Row */}
            <View style={styles.leaderboardHeaderRow}>
              <Text style={[styles.leaderboardHeaderCell, { flex: 1, textAlign: 'center' }]}>#</Text>
              <Text style={[styles.leaderboardHeaderCell, { flex: 4 }]}>Hráč</Text>
              <Text style={[styles.leaderboardHeaderCell, { flex: 2, textAlign: 'center' }]}>Kills</Text>
              <Text style={[styles.leaderboardHeaderCell, { flex: 2, textAlign: 'center' }]}>Deaths</Text>
            </View>
            
            {/* Data Rows */}
            {MOCK_LEADERBOARD.map((item, index) => (
              <View key={item.id} style={[
                styles.leaderboardRow, 
                index !== MOCK_LEADERBOARD.length - 1 && styles.leaderboardRowBorder
              ]}>
                <Text style={[styles.leaderboardCell, { flex: 1, textAlign: 'center', color: item.rank === 1 ? COLORS.accentGreen : COLORS.textSecondary }]}>{item.rank}</Text>
                <Text style={[styles.leaderboardCell, { flex: 4, color: item.id === currentUser.id ? COLORS.textPrimary : COLORS.textSecondary, fontFamily: item.id === currentUser.id ? 'Goldman_700Bold' : 'Goldman_400Regular' }]}>{item.name}</Text>
                <Text style={[styles.leaderboardCell, { flex: 2, textAlign: 'center', color: COLORS.accentGreen }]}>{item.kills}</Text>
                <Text style={[styles.leaderboardCell, { flex: 2, textAlign: 'center', color: COLORS.accentRed }]}>{item.deaths}</Text>
              </View>
            ))}
          </View>
        </View>

        {/* HISTORY SECTION */}
        <View style={styles.sectionContainer}>
          <Text style={styles.sectionTitle}>PŘEDCHOZÍ HRY</Text>
          {MOCK_HISTORY.map((game) => (
            <View key={game.id} style={styles.historyCard}>
              <View style={styles.historyHeader}>
                <Text style={styles.historyMode}>{game.mode}</Text>
                <Text style={styles.historyDate}>{game.date}</Text>
              </View>
              <Text style={styles.historyMap}>{game.map}</Text>
              <View style={styles.historyStats}>
                 <View style={styles.historyStatBadge}>
                    <Text style={[styles.historyStatLabel, { color: COLORS.textSecondary }]}>RANK</Text>
                    <Text style={[styles.historyStatValue, { color: game.rank === 1 ? COLORS.accentGreen : COLORS.textPrimary }]}>#{game.rank}</Text>
                 </View>
                 <View style={styles.historyStatBadge}>
                    <Text style={[styles.historyStatLabel, { color: COLORS.accentGreen }]}>KILLS</Text>
                    <Text style={styles.historyStatValue}>{game.kills}</Text>
                 </View>
                 <View style={styles.historyStatBadge}>
                    <Text style={[styles.historyStatLabel, { color: COLORS.accentRed }]}>DEATHS</Text>
                    <Text style={styles.historyStatValue}>{game.deaths}</Text>
                 </View>
              </View>
            </View>
          ))}
        </View>

        {/* ACTION BUTTON */}
        <TouchableOpacity style={styles.actionButton} onPress={() => Alert.alert('Hledání hry', 'Připojování k serveru...')}>
          <Text style={styles.actionButtonText}>OPAKOVAT</Text>
        </TouchableOpacity>

      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    backgroundColor: COLORS.background,
    alignItems: 'center',
    justifyContent: 'center',
  },
  container: {
    flex: 1,
    backgroundColor: COLORS.background,
    paddingTop: Platform.OS === 'android' ? 30 : 0,
  },
  scrollContent: {
    padding: 20,
    paddingBottom: 40,
  },
  header: {
    marginTop: 20,
    marginBottom: 30,
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: 28, // Smaller title
    fontFamily: 'Goldman_700Bold',
    color: COLORS.textPrimary,
    letterSpacing: 2,
    marginBottom: 5,
    textAlign: 'center',
  },
  gameInfoContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 5,
  },
  gameInfoText: {
    color: COLORS.accentGreen,
    fontFamily: 'Goldman_700Bold',
    fontSize: 14,
    letterSpacing: 1,
  },
  separator: {
    color: COLORS.textSecondary,
    marginHorizontal: 10,
  },
  dateText: {
    fontSize: 12,
    fontFamily: 'Goldman_400Regular',
    color: COLORS.textSecondary,
  },
  
  // Player Card
  playerCard: {
    backgroundColor: COLORS.cardBackground,
    borderRadius: 12,
    padding: 20,
    marginBottom: 20,
    borderWidth: 1,
    borderColor: COLORS.border,
    alignItems: 'center',
    width: '100%',
  },
  playerLabel: {
    color: COLORS.textSecondary,
    fontSize: 10,
    fontFamily: 'Goldman_400Regular',
    letterSpacing: 2,
    marginBottom: 5,
    textTransform: 'uppercase',
  },
  playerName: {
    color: COLORS.textPrimary,
    fontSize: 32,
    fontFamily: 'Goldman_700Bold',
    marginBottom: 10,
  },
  rankBadge: {
    backgroundColor: COLORS.background,
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: COLORS.accentGreen,
  },
  rankText: {
    color: COLORS.accentGreen,
    fontFamily: 'Goldman_700Bold',
    fontSize: 12,
  },

  // Grid Styles
  gridContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    marginBottom: 20,
  },
  gridItem: {
    width: '48%', // Approx half with spacing
    backgroundColor: COLORS.cardBackground,
    borderRadius: 12,
    paddingVertical: 20,
    paddingHorizontal: 10,
    marginBottom: 15,
    borderWidth: 1,
    borderColor: COLORS.border,
    alignItems: 'center',
  },
  gridLabel: {
    color: COLORS.textSecondary,
    fontSize: 10,
    fontFamily: 'Goldman_400Regular',
    letterSpacing: 1,
    marginBottom: 5,
    textTransform: 'uppercase',
    textAlign: 'center',
  },
  gridValue: {
    color: COLORS.textPrimary,
    fontSize: 32,
    fontFamily: 'Goldman_700Bold',
  },
  gridSubtext: {
    color: COLORS.textSecondary,
    fontSize: 10,
    marginTop: 5,
    fontFamily: 'Goldman_400Regular',
  },

  // Leaderboard Styles
  sectionContainer: {
    marginBottom: 30,
  },
  sectionTitle: {
    color: COLORS.textSecondary,
    fontSize: 14,
    fontFamily: 'Goldman_400Regular',
    letterSpacing: 2,
    marginBottom: 10,
    textTransform: 'uppercase',
  },
  leaderboardContainer: {
    backgroundColor: COLORS.cardBackground, // Slightly lighter than pure black
    borderRadius: 12,
    borderWidth: 1,
    borderColor: COLORS.border,
    overflow: 'hidden',
  },
  
  // History Styles
  historyCard: {
    backgroundColor: COLORS.cardBackground,
    borderRadius: 8,
    padding: 15,
    marginBottom: 10,
    borderLeftWidth: 4,
    borderLeftColor: COLORS.accentGreen, // Default accent
  },
  historyHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 5,
  },
  historyMode: {
    color: COLORS.textPrimary,
    fontFamily: 'Goldman_700Bold',
    fontSize: 14,
  },
  historyDate: {
    color: COLORS.textSecondary,
    fontFamily: 'Goldman_400Regular',
    fontSize: 12,
  },
  historyMap: {
    color: COLORS.textSecondary,
    fontFamily: 'Goldman_400Regular',
    fontSize: 10,
    textTransform: 'uppercase',
    marginBottom: 10,
  },
  historyStats: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    borderTopWidth: 1,
    borderTopColor: COLORS.border,
    paddingTop: 10,
  },
  historyStatBadge: {
    alignItems: 'center',
    flex: 1,
  },
  historyStatLabel: {
    fontSize: 10,
    fontFamily: 'Goldman_400Regular',
    marginBottom: 2,
  },
  historyStatValue: {
    color: COLORS.textPrimary,
    fontSize: 16,
    fontFamily: 'Goldman_700Bold',
  },
  leaderboardHeaderRow: {
    flexDirection: 'row',
    backgroundColor: '#09090b', // Zinc-950
    paddingVertical: 12,
    paddingHorizontal: 15,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
  },
  leaderboardHeaderCell: {
    color: COLORS.textSecondary,
    fontSize: 10,
    fontFamily: 'Goldman_400Regular',
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  leaderboardRow: {
    flexDirection: 'row',
    paddingVertical: 12,
    paddingHorizontal: 15,
    alignItems: 'center',
  },
  leaderboardRowBorder: {
    borderBottomWidth: 1,
    borderBottomColor: '#27272a', // Zinc-800
  },
  leaderboardCell: {
    color: COLORS.textPrimary,
    fontSize: 14,
    fontFamily: 'Goldman_400Regular',
  },

  // Action Button
  actionButton: {
    backgroundColor: COLORS.accentGreen,
    borderRadius: 8,
    paddingVertical: 16,
    alignItems: 'center',
    shadowColor: COLORS.accentGreen,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.5,
    shadowRadius: 10,
    elevation: 5,
    marginBottom: 20
  },
  actionButtonText: {
    color: '#000',
    fontSize: 18,
    fontFamily: 'Goldman_700Bold',
    letterSpacing: 2,
    textTransform: 'uppercase',
  },

  // Login Screen Styles
  loginContainer: {
    flex: 1,
    backgroundColor: COLORS.background,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  loginLogo: {
    fontSize: 48,
    fontFamily: 'Goldman_700Bold',
    color: COLORS.accentGreen,
    marginBottom: 10,
    textAlign: 'center',
    textShadowColor: COLORS.accentGreen,
    textShadowOffset: { width: 0, height: 0 },
    textShadowRadius: 10,
  },
  loginSubtitle: {
    fontSize: 16,
    fontFamily: 'Goldman_400Regular',
    color: COLORS.textPrimary,
    marginBottom: 50,
    letterSpacing: 2,
  },
  inputContainer: {
    width: '100%',
    marginBottom: 30,
  },
  inputLabel: {
    color: COLORS.textSecondary,
    fontFamily: 'Goldman_700Bold',
    fontSize: 12,
    marginBottom: 10,
    letterSpacing: 1,
  },
  input: {
    backgroundColor: COLORS.cardBackground,
    borderWidth: 1,
    borderColor: COLORS.border,
    borderRadius: 8,
    padding: 15,
    color: COLORS.textPrimary,
    fontFamily: 'Goldman_700Bold',
    fontSize: 24,
    textAlign: 'center',
    letterSpacing: 5,
  },
  loginButton: {
    backgroundColor: COLORS.accentGreen,
    width: '100%',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
    marginBottom: 20,
  },
  loginFooter: {
    color: COLORS.textSecondary,
    fontSize: 10,
    fontFamily: 'Goldman_400Regular',
    textAlign: 'center',
    opacity: 0.5,
  },

  // Top Bar & Navigation
  topBar: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    marginBottom: 10,
  },
  logoutButton: {
    padding: 8,
  },
  logoutText: {
    color: COLORS.accentRed,
    fontFamily: 'Goldman_700Bold',
    fontSize: 12,
  },
});
