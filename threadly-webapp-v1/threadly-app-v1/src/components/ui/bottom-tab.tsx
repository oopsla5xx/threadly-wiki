import type { BottomTabBarProps } from "expo-router/js-tabs"
import { Heart, House, Plus, Send, UserRound } from "lucide-react-native"
import { THEME } from "@/lib/theme"
import { useEffect, useRef, useState } from "react"
import { Pressable, View, type LayoutChangeEvent } from "react-native"
import Animated, {
  useAnimatedStyle,
  useSharedValue,
  withSequence,
  withSpring,
  withTiming,
} from "react-native-reanimated"
import { useSafeAreaInsets } from "react-native-safe-area-context"

const BAR_PADDING_X = 20
const INDICATOR_INSET = 4

const iconByRoute = {
  home: House,
  message: Send,
  new_thread: Plus,
  activity: Heart,
  profile: UserRound,
} as const

export function BottomTabBar({ state, descriptors, navigation }: BottomTabBarProps) {
  const insets = useSafeAreaInsets()
  const [tabBarWidth, setTabBarWidth] = useState(0)
  const indicatorOffset = useSharedValue(0)
  const indicatorScale = useSharedValue(1)
  const previousIndex = useRef(state.index)
  const activeIndex = state.index
  const slotWidth = Math.max(0, tabBarWidth - BAR_PADDING_X * 2) / state.routes.length
  const indicatorWidth = Math.max(0, slotWidth - INDICATOR_INSET * 2)

  useEffect(() => {
    // Width-only changes (first layout, rotation) must not animate or bounce.
    if (previousIndex.current === activeIndex) {
      indicatorOffset.value = activeIndex * slotWidth
      return
    }

    previousIndex.current = activeIndex
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
      className="mx-4 h-20 flex-row items-center justify-center rounded-[2.4rem] bg-white shadow-2xl"
      style={{ paddingHorizontal: BAR_PADDING_X, marginBottom: Math.max(insets.bottom, 20) }}
      onLayout={handleLayout}
    >
      {indicatorWidth > 0 ? (
        <Animated.View
          pointerEvents="none"
          className="absolute h-14 rounded-full border border-slate-100 bg-slate-100"
          style={[{ left: BAR_PADDING_X + INDICATOR_INSET, width: indicatorWidth }, indicatorStyle]}
        />
      ) : null}

      {state.routes.map((route, routeIndex) => {
        const Icon = iconByRoute[route.name as keyof typeof iconByRoute]
        if (!Icon) return null

        const isFocused = state.index === routeIndex
        const label = descriptors[route.key].options.title ?? route.name

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
              color={THEME.light.foreground}
              fill={isFocused && route.name !== "new_thread" ? THEME.light.foreground : "none"}
              size={route.name === "new_thread" ? 34 : 29}
              strokeWidth={2.2}
            />
          </Pressable>
        )
      })}
    </View>
  )
}
