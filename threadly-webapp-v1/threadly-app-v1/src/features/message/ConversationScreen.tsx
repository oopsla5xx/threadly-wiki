import { Text } from "@/components/ui/text"
import { View } from "react-native"

type ConversationScreenProps = {
  id: string
}

export function ConversationScreen({ id }: ConversationScreenProps) {
  return (
    <View className="flex-1 items-center justify-center bg-white px-6">
      <Text variant="h3">Conversation {id}</Text>
    </View>
  )
}
