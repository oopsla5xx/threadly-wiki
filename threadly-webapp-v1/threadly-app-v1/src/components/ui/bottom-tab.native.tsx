import type { BottomTabBarProps } from "expo-router/js-tabs"
import { Heart, House, Plus, Send, UserRound } from "lucide-react-native"
import { useEffect, useState } from "react"
import { Pressable, View, type LayoutChangeEvent } from "react-native"
import Animated, {
  useAnimatedStyle,
  useSharedValue,
  withSequence,
  withSpring,
  withTiming,
} from "react-native-reanimated"
import { useSafeAreaInsets } from "react-native-safe-area-context"

const tabRoutes = ["home", "message", "new_thread", "activity", "profile"] as const

const iconByRoute = {
  home: House,
  message: Send,
  new_thread: Plus,
  activity: Heart,
  profile: UserRound,
} as const

export function HomeBottomTabs({ state, descriptors, navigation }: BottomTabBarProps) {
  const insets = useSafeAreaInsets()
  const [tabBarWidth, setTabBarWidth] = useState(0)
  const indicatorOffset = useSharedValue(0)
  const indicatorScale = useSharedValue(1)
  const activeRouteName = state.routes[state.index]?.name
  const activeIndex = Math.max(0, tabRoutes.indexOf(activeRouteName as (typeof tabRoutes)[number]))
  const slotWidth = Math.max(0, tabBarWidth - 40) / tabRoutes.length
  const indicatorWidth = Math.max(0, slotWidth - 8)

  useEffect(() => {
    indicatorOffset.value = withSpring(activeIndex * slotWidth, {
      damping: 18,
      stiffness: 180,
      mass: 0.8,
    })
    indicatorScale.value = withSequence(
      withTiming(1.14, { duration: 100 }),
      withSpring(1, { damping: 16, stiffness: 220 }),
    )
  }, [activeIndex, indicatorOffset, indicatorScale, slotWidth])

  const indicatorStyle = useAnimatedStyle(() => ({
    transform: [{ translateX: indicatorOffset.value }, { scale: indicatorScale.value }],
  }))

  const handleLayout = (event: LayoutChangeEvent) => {
    setTabBarWidth(event.nativeEvent.layout.width)
  }

  return (
    <View
      className="mx-4 h-20 flex-row items-center justify-center rounded-[2.4rem] bg-white px-5 shadow-2xl"
      style={{ marginBottom: Math.max(insets.bottom, 20) }}
      onLayout={handleLayout}
    >
      {indicatorWidth > 0 ? (
        <Animated.View
          pointerEvents="none"
          className="absolute left-2 h-14 rounded-b-[2rem] rounded-l-[4rem] rounded-r-[4rem] rounded-t-[2rem] border border-slate-100 bg-slate-100 px-11 py-4"
          style={[{ width: indicatorWidth }, indicatorStyle]}
        />
      ) : null}

      {tabRoutes.map((routeName) => {
        const route = state.routes.find((item) => item.name === routeName)
        if (!route) return null

        const routeIndex = state.routes.findIndex((item) => item.key === route.key)
        const isFocused = state.index === routeIndex
        const label = descriptors[route.key].options.title ?? route.name
        const Icon = iconByRoute[routeName]

        return (
          <Pressable
            key={route.key}
            accessibilityLabel={label}
            accessibilityRole="tab"
            accessibilityState={{ selected: isFocused }}
            className="h-14 flex-1 items-center justify-center active:opacity-50"
            onPress={() => {
              const event = navigation.emit({
                type: "tabPress",
                target: route.key,
                canPreventDefault: true,
              })

              if (!isFocused && !event.defaultPrevented) {
                navigation.navigate(route.name, route.params)
              }
            }}
            onLongPress={() => navigation.emit({ type: "tabLongPress", target: route.key })}
          >
            <Icon
              color="#0a0a0a"
              fill={
                isFocused &&
                (routeName === "home" ||
                  routeName === "activity" ||
                  routeName === "message" ||
                  routeName === "profile")
                  ? "#0a0a0a"
                  : "none"
              }
              size={routeName === "new_thread" ? 34 : 29}
              strokeWidth={2.2}
            />
          </Pressable>
        )
      })}
    </View>
  )
}
