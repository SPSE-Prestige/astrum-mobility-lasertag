import { SafeAreaView, Text, View } from "react-native";
import { StatusBar } from "expo-status-bar";

export default function Home() {
  return (
    <SafeAreaView className="flex-1 bg-[#020303]">
      <StatusBar barStyle="light-content" backgroundColor="#020303" />
      <View className="flex-1 justify-center items-center px-4">
        <Text className="text-white text-4xl font-bold mb-4">LaserTag</Text>
        <Text className="text-gray-400 text-lg text-center">
          Mobilní aplikace pro race control systém
        </Text>
      </View>
    </SafeAreaView>
  );
}
