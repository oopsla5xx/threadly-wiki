import { Text } from "@/components/ui/text"
import { View } from "react-native"

export function NewThreadScreen() {
  return (
    <View className="flex-1 items-center justify-center bg-white px-6">
      <Text variant="h3">Create a thread</Text>
      <Text className="mt-2 text-center text-slate-500">Share what is on your mind.</Text>
    </View>
  )
}
