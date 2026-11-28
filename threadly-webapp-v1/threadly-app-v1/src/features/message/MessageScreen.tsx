import { Text } from "@/components/ui/text"
import { Link } from "expo-router"
import { View } from "react-native"

export function MessageScreen() {
  return (
    <View className="flex-1 items-center justify-center bg-white px-6">
      <Text variant="h3">Messages</Text>
      <Text className="mt-2 text-center text-slate-500">Your conversations will appear here.</Text>
      <Link href="/message/123" className="mt-4 text-blue-500">
        Open conversation
      </Link>
    </View>
  )
}
