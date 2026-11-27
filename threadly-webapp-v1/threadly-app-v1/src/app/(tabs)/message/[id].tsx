import { ConversationScreen } from "@/features/message/ConversationScreen"
import { useLocalSearchParams } from "expo-router"

export default function ConversationRoute() {
  const { id } = useLocalSearchParams<{ id: string }>()

  return <ConversationScreen id={id} />
}
