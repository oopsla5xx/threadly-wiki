import { Text } from "@/components/ui/text"
import { View } from "react-native"

export function HomeScreen() {
  return (
    <View className="flex-1 items-center justify-center bg-white">
      <Text className="text-xl font-bold text-slate-900">Threadly</Text>
      <Text className="mt-2 text-slate-600">Your app is ready.</Text>
    </View>
  )
}
