import { Redirect, Tabs } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Platform, Pressable, StyleSheet, View } from 'react-native';
import { BlurView } from 'expo-blur';
import * as Haptics from 'expo-haptics';
import { useTheme } from '../../src/components/ui/theme-provider';
import { useAuth } from '../../src/lib/AuthContext';
import {
  spacing,
  sizes,
  sharedColors,
  androidGlassBg,
  tabActiveBg,
  tabInactiveTint,
} from '../../src/components/ui/tokens';

const TAB_BAR_HEIGHT = sizes[13]; // 52 — pill height
const TAB_ACTIVE_INSET = sizes[1]; // 4 — equal gap on all sides of the active highlight
const TAB_ACTIVE_SIZE = TAB_BAR_HEIGHT - TAB_ACTIVE_INSET * 2; // 44 — highlight fills pill minus inset
const BUTTON_BORDER_RADIUS = TAB_ACTIVE_SIZE / 2;

const styles = StyleSheet.create({
  blurView: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    borderRadius: TAB_BAR_HEIGHT / 2,
    overflow: 'hidden',
  },
  tabButtonOuter: {
    flex: 1,
    alignSelf: 'stretch',
    alignItems: 'center',
    justifyContent: 'center',
  },
  tabButtonInner: {
    width: TAB_ACTIVE_SIZE,
    height: TAB_ACTIVE_SIZE,
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: BUTTON_BORDER_RADIUS,
  },
  addButtonOuter: {
    flex: 1,
    alignSelf: 'stretch',
    alignItems: 'center',
    justifyContent: 'center',
  },
  addButtonInner: {
    width: TAB_ACTIVE_SIZE,
    height: TAB_ACTIVE_SIZE,
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: BUTTON_BORDER_RADIUS,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.35,
    shadowRadius: 8,
    elevation: 5,
  },
});

function TabBarBackground() {
  const { resolved } = useTheme();
  if (Platform.OS === 'android') {
    return null;
  }
  return (
    <BlurView
      intensity={80}
      tint={resolved === 'dark' ? 'dark' : 'light'}
      style={styles.blurView}
    />
  );
}

// Completely replaces the internal tab button so we own the full layout.
// This bypasses the library's tabVerticalUiKit style (justifyContent: 'flex-start')
// which was preventing icons from centering.
function makeTabButton(
  iconName: keyof typeof Ionicons.glyphMap,
  focusedIconName: keyof typeof Ionicons.glyphMap,
  hapticStyle: Haptics.ImpactFeedbackStyle,
  resolved: 'light' | 'dark',
  activeTintColor: string,
  inactiveTintColor: string
) {
  const tintBg = resolved === 'dark' ? tabActiveBg.dark : tabActiveBg.light;
  return function TabButton({ onPress, onLongPress, children: _children, ...rest }: any) {
    const focused = rest['aria-selected'] as boolean;
    const color = focused ? activeTintColor : inactiveTintColor;
    const bgColor = focused ? tintBg : sharedColors.input;
    return (
      <Pressable
        onPress={() => {
          Haptics.impactAsync(hapticStyle);
          onPress?.();
        }}
        onLongPress={onLongPress}
        style={styles.tabButtonOuter}
      >
        <View style={[styles.tabButtonInner, { backgroundColor: bgColor }]}>
          <Ionicons name={focused ? focusedIconName : iconName} size={24} color={color} />
        </View>
      </Pressable>
    );
  };
}

function makeAddButton(primaryColor: string, hapticStyle: Haptics.ImpactFeedbackStyle) {
  return function AddButton({ onPress, onLongPress }: any) {
    return (
      <Pressable
        onPress={() => {
          Haptics.impactAsync(hapticStyle);
          onPress?.();
        }}
        onLongPress={onLongPress}
        style={styles.addButtonOuter}
      >
        <View
          style={[
            styles.addButtonInner,
            { backgroundColor: primaryColor, shadowColor: primaryColor },
          ]}
        >
          <Ionicons name="add" size={26} color={sharedColors.primaryForeground} />
        </View>
      </Pressable>
    );
  };
}

export default function TabLayout() {
  const { colors, resolved } = useTheme();
  const { hasToken, isChecking } = useAuth();

  if (!isChecking && !hasToken) {
    return <Redirect href="/auth" />;
  }

  const androidBg = resolved === 'dark' ? androidGlassBg.dark : androidGlassBg.light;
  const active = colors.primary;
  const inactive = resolved === 'dark' ? tabInactiveTint.dark : tabInactiveTint.light;
  const tabBarSideMargin = spacing[4]; // 16 — matches filter bar horizontal inset

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: active,
        tabBarInactiveTintColor: inactive,
        tabBarShowLabel: false,
        tabBarStyle: {
          position: 'absolute',
          bottom: spacing[6],
          left: 0,
          right: 0,
          marginHorizontal: tabBarSideMargin,
          height: TAB_BAR_HEIGHT,
          // The pill floats above the home indicator, so cancel the library's
          // automatic paddingBottom (= insets.bottom) which would shrink the
          // content area and push icons to the top.
          paddingBottom: 0,
          borderRadius: TAB_BAR_HEIGHT / 2,
          backgroundColor: Platform.OS === 'android' ? androidBg : sharedColors.input,
          borderTopWidth: 0,
          elevation: Platform.OS === 'android' ? 8 : 0,
          shadowColor: colors.foreground,
          shadowOffset: { width: 0, height: 8 },
          shadowOpacity: resolved === 'dark' ? 0.4 : 0.15,
          shadowRadius: 24,
          overflow: Platform.OS === 'android' ? 'hidden' : 'visible',
        },
        tabBarBackground: () => <TabBarBackground />,
        headerShown: false,
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
          tabBarButton: makeTabButton(
            'home-outline',
            'home',
            Haptics.ImpactFeedbackStyle.Light,
            resolved,
            active,
            inactive
          ),
        }}
      />
      <Tabs.Screen
        name="history"
        options={{
          title: 'History',
          tabBarButton: makeTabButton(
            'time-outline',
            'time',
            Haptics.ImpactFeedbackStyle.Light,
            resolved,
            active,
            inactive
          ),
        }}
      />
      <Tabs.Screen
        name="add"
        options={{
          title: 'Add',
          tabBarButton: makeAddButton(colors.primary, Haptics.ImpactFeedbackStyle.Medium),
        }}
      />
      <Tabs.Screen
        name="friends"
        options={{
          title: 'Friends',
          tabBarButton: makeTabButton(
            'people-outline',
            'people',
            Haptics.ImpactFeedbackStyle.Light,
            resolved,
            active,
            inactive
          ),
        }}
      />
      <Tabs.Screen
        name="account"
        options={{
          title: 'Account',
          tabBarButton: makeTabButton(
            'person-outline',
            'person',
            Haptics.ImpactFeedbackStyle.Light,
            resolved,
            active,
            inactive
          ),
        }}
      />
    </Tabs>
  );
}
