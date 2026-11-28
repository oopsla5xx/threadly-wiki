import { Text } from "@/components/ui/text"
import { View } from "react-native"

export function ProfileScreen() {
  return (
    <View className="flex-1 items-center justify-center bg-white px-6">
      <Text variant="h3">Profile</Text>
      <Text className="mt-2 text-center text-slate-500">Manage your profile and account.</Text>
    </View>
  )
}
