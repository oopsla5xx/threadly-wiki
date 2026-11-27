import { Text } from "@/components/ui/text"
import { Link } from "expo-router"
import { View } from "react-native"

export default function NotFoundScreen() {
  return (
    <View className="flex-1 items-center justify-center gap-4 bg-white px-6">
      <Text variant="h2">Screen not found</Text>
      <Link href="/(tabs)/home" className="font-medium text-slate-900 underline">
        Go home
      </Link>
    </View>
  )
}
