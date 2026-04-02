import { useState, useRef, useCallback } from 'react';
import {
  SafeAreaView,
  Text,
  View,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  Keyboard,
} from 'react-native';
import { useRouter } from 'expo-router';
import { getPlayerSession } from '../lib/api';

const VALID_CHARS = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789';

export default function LoginScreen() {
  const router = useRouter();
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const inputRef = useRef<TextInput>(null);

  const handleCodeChange = useCallback(
    (text: string) => {
      const filtered = text
        .toUpperCase()
        .split('')
        .filter((c) => VALID_CHARS.includes(c))
        .join('')
        .slice(0, 6);

      setCode(filtered);
      setError('');

      if (filtered.length === 6) {
        Keyboard.dismiss();
        submitCode(filtered);
      }
    },
    []
  );

  const submitCode = async (pin: string) => {
    if (pin.length !== 6) return;
    setLoading(true);
    setError('');
    try {
      await getPlayerSession(pin);
      router.push({ pathname: '/dashboard', params: { code: pin } });
    } catch {
      setError('Neplatný kód');
    } finally {
      setLoading(false);
    }
  };

  const boxes = Array.from({ length: 6 }, (_, i) => code[i] || '');

  return (
    <SafeAreaView className="flex-1 bg-[#020303]">
      <View className="flex-1 justify-center items-center px-6">
        {/* Logo */}
        <Text className="text-6xl mb-2">⊕</Text>
        <Text className="text-white text-4xl font-bold tracking-widest mb-1">
          LASER TAG
        </Text>
        <Text className="text-[#00F0FF] text-base mb-10 opacity-70">
          Zadej svůj kód hráče
        </Text>

        {/* PIN boxes */}
        <TouchableOpacity
          activeOpacity={1}
          onPress={() => inputRef.current?.focus()}
          className="flex-row justify-center mb-6"
        >
          {boxes.map((char, i) => (
            <View
              key={i}
              className={`w-12 h-14 mx-1 rounded-lg border-2 justify-center items-center ${
                error
                  ? 'border-[#FF3B30]'
                  : i === code.length
                    ? 'border-[#00F0FF]'
                    : char
                      ? 'border-[#00F0FF]/50'
                      : 'border-[#1a1a2e]'
              } bg-[#0A0F0F]`}
            >
              <Text className="text-white text-2xl font-bold">{char}</Text>
            </View>
          ))}
        </TouchableOpacity>

        {/* Hidden input */}
        <TextInput
          ref={inputRef}
          value={code}
          onChangeText={handleCodeChange}
          autoCapitalize="characters"
          autoCorrect={false}
          maxLength={6}
          className="absolute opacity-0 h-0 w-0"
          autoFocus
        />

        {/* Error */}
        {error ? (
          <Text className="text-[#FF3B30] text-sm mb-4 font-semibold">
            {error}
          </Text>
        ) : (
          <View className="h-5 mb-4" />
        )}

        {/* Submit button */}
        <TouchableOpacity
          onPress={() => submitCode(code)}
          disabled={code.length !== 6 || loading}
          className={`w-full py-4 rounded-xl items-center ${
            code.length === 6 && !loading
              ? 'bg-[#00F0FF]'
              : 'bg-[#00F0FF]/20'
          }`}
        >
          {loading ? (
            <ActivityIndicator color="#020303" />
          ) : (
            <Text
              className={`text-lg font-bold tracking-wider ${
                code.length === 6 ? 'text-[#020303]' : 'text-[#00F0FF]/40'
              }`}
            >
              PŘIPOJIT SE
            </Text>
          )}
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}
