import { Text } from "@/components/ui/text"
import { View } from "react-native"

export function ActivityScreen() {
  return (
    <View className="flex-1 items-center justify-center bg-white px-6">
      <Text variant="h3">Activity</Text>
      <Text className="mt-2 text-center text-slate-500">
        Your latest activity will appear here.
      </Text>
    </View>
  )
}
