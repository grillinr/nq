import { Redirect, Tabs } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Platform, Pressable, View } from 'react-native';
import { BlurView } from 'expo-blur';
import * as Haptics from 'expo-haptics';
import { useTheme } from '../../src/components/ui/theme-provider';
import { useAuth } from '../../src/lib/AuthContext';

function TabBarBackground() {
  const { resolved } = useTheme();
  if (Platform.OS === 'android') {
    return null;
  }
  return (
    <BlurView
      intensity={80}
      tint={resolved === 'dark' ? 'dark' : 'light'}
      style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        borderRadius: 32,
        overflow: 'hidden',
      }}
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
  const tintBg = resolved === 'dark' ? 'rgba(38,99,217,0.25)' : 'rgba(38,99,217,0.12)';
  return function TabButton({ onPress, onLongPress, children: _children, ...rest }: any) {
    const focused = rest['aria-selected'] as boolean;
    const color = focused ? activeTintColor : inactiveTintColor;
    return (
      <Pressable
        onPress={() => {
          Haptics.impactAsync(hapticStyle);
          onPress?.();
        }}
        onLongPress={onLongPress}
        style={{ flex: 1, alignSelf: 'stretch', alignItems: 'center', justifyContent: 'center' }}
      >
        <View
          style={{
            width: 44,
            height: 44,
            alignItems: 'center',
            justifyContent: 'center',
            borderRadius: 22,
            backgroundColor: focused ? tintBg : 'transparent',
          }}
        >
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
        style={{ flex: 1, alignSelf: 'stretch', alignItems: 'center', justifyContent: 'center' }}
      >
        <View
          style={{
            width: 44,
            height: 44,
            alignItems: 'center',
            justifyContent: 'center',
            borderRadius: 22,
            backgroundColor: primaryColor,
            shadowColor: primaryColor,
            shadowOffset: { width: 0, height: 4 },
            shadowOpacity: 0.35,
            shadowRadius: 8,
            elevation: 5,
          }}
        >
          <Ionicons name="add" size={26} color="#ffffff" />
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

  const androidBg = resolved === 'dark' ? 'rgba(28,28,30,0.95)' : 'rgba(255,255,255,0.95)';
  const active = colors.primary;
  const inactive = resolved === 'dark' ? 'rgba(255,255,255,0.4)' : 'rgba(0,0,0,0.3)';

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: active,
        tabBarInactiveTintColor: inactive,
        tabBarShowLabel: false,
        tabBarStyle: {
          position: 'absolute',
          bottom: 24,
          left: 24,
          right: 24,
          height: 64,
          // The pill floats above the home indicator, so cancel the library's
          // automatic paddingBottom (= insets.bottom) which would shrink the
          // content area and push icons to the top.
          paddingBottom: 0,
          borderRadius: 32,
          backgroundColor: Platform.OS === 'android' ? androidBg : 'transparent',
          borderTopWidth: 0,
          elevation: Platform.OS === 'android' ? 8 : 0,
          shadowColor: '#000',
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
