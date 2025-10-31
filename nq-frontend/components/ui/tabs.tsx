import React, { createContext, useContext } from 'react';
import { View, Pressable, StyleSheet } from 'react-native';
import { colors, spacing, radii } from './tokens';

interface TabsContextType {
  value: string;
  onValueChange: (value: string) => void;
}

const TabsContext = createContext<TabsContextType | null>(null);

interface TabsProps {
  value: string;
  onValueChange: (value: string) => void;
  children: React.ReactNode;
}

interface TabsListProps {
  children: React.ReactNode;
}

interface TabsTriggerProps {
  value: string;
  children: React.ReactNode;
}

interface TabsContentProps {
  value: string;
  children: React.ReactNode;
}

export function Tabs({ value, onValueChange, children }: TabsProps) {
  return (
    <TabsContext.Provider value={{ value, onValueChange }}>
      <View style={styles.tabs}>{children}</View>
    </TabsContext.Provider>
  );
}

export function TabsList({ children }: TabsListProps) {
  return <View style={styles.tabsList}>{children}</View>;
}

export function TabsTrigger({ value, children }: TabsTriggerProps) {
  const context = useContext(TabsContext);
  if (!context) throw new Error('TabsTrigger must be used within Tabs');

  const { value: activeValue, onValueChange } = context;
  const isActive = activeValue === value;

  return (
    <Pressable
      style={[styles.tabsTrigger, isActive && styles.tabsTriggerActive]}
      onPress={() => onValueChange(value)}
    >
      {children}
    </Pressable>
  );
}

export function TabsContent({ value, children }: TabsContentProps) {
  const context = useContext(TabsContext);
  if (!context) throw new Error('TabsContent must be used within Tabs');

  const { value: activeValue } = context;
  if (activeValue !== value) return null;

  return <View style={styles.tabsContent}>{children}</View>;
}

const styles = StyleSheet.create({
  tabs: {
    // Container for tabs
  },
  tabsList: {
    flexDirection: 'row',
    backgroundColor: colors.muted,
    borderRadius: radii.md,
    padding: spacing[1],
  },
  tabsTrigger: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: spacing[2],
    paddingHorizontal: spacing[3],
    borderRadius: radii.sm,
  },
  tabsTriggerActive: {
    backgroundColor: colors.background,
    shadowColor: colors.foreground,
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  tabsContent: {
    // Content container
  },
});