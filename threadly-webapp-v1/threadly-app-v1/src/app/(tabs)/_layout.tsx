import { BottomTabBar } from "@/components/ui/bottom-tab"
import { Tabs } from "expo-router"

export default function TabsLayout() {
  return (
    <Tabs
      initialRouteName="home"
      screenOptions={{ headerShown: false }}
      tabBar={(props) => <BottomTabBar {...props} />}
    >
      <Tabs.Screen name="home" options={{ title: "Home" }} />
      <Tabs.Screen name="message" options={{ title: "Messages" }} />
      <Tabs.Screen name="new_thread" options={{ title: "New thread" }} />
      <Tabs.Screen name="activity" options={{ title: "Activity" }} />
      <Tabs.Screen name="profile" options={{ title: "Profile" }} />
    </Tabs>
  )
}
