import React, { useMemo, useState } from 'react';
import { View, Text, Pressable, StyleSheet, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Modal } from './modal';
import { spacing, radii, fontSize, fontWeights, ColorPalette } from './tokens';
import { useTheme } from './theme-provider';

interface SelectProps {
  value?: string;
  onValueChange: (value: string) => void;
  children: React.ReactNode;
}

interface SelectTriggerProps {
  children: React.ReactNode;
}

interface SelectValueProps {
  placeholder?: string;
}

interface SelectContentProps {
  children: React.ReactNode;
}

interface SelectItemProps {
  value: string;
  children: React.ReactNode;
}

export function Select({ value, onValueChange, children }: SelectProps) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <View style={selectContainerStyle}>
      {React.Children.map(children, (child) =>
        React.isValidElement(child)
          ? React.cloneElement(child as React.ReactElement<any>, {
              isOpen,
              setIsOpen,
              value,
              onValueChange,
            })
          : child
      )}
    </View>
  );
}

const selectContainerStyle = { position: 'relative' as const };

function createTriggerStyles(colors: ColorPalette) {
  return StyleSheet.create({
    selectTrigger: {
      flexDirection: 'row',
      alignItems: 'center',
      justifyContent: 'space-between',
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[2],
      borderWidth: 1,
      borderColor: colors.border,
      borderRadius: radii.md,
      backgroundColor: colors.background,
    },
  });
}

export function SelectTrigger({ children, ...props }: SelectTriggerProps & any) {
  const { isOpen, setIsOpen } = props;
  const { colors } = useTheme();
  const styles = useMemo(() => createTriggerStyles(colors), [colors]);

  return (
    <Pressable style={styles.selectTrigger} onPress={() => setIsOpen(!isOpen)}>
      {children}
      <Ionicons
        name={isOpen ? 'chevron-up' : 'chevron-down'}
        size={16}
        color={colors.mutedForeground}
      />
    </Pressable>
  );
}

function createValueStyles(colors: ColorPalette) {
  return StyleSheet.create({
    selectValue: { fontSize: fontSize.base, color: colors.foreground },
  });
}

export function SelectValue({ placeholder }: SelectValueProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createValueStyles(colors), [colors]);
  return <Text style={styles.selectValue}>{placeholder}</Text>;
}

function createContentStyles(colors: ColorPalette) {
  return StyleSheet.create({
    selectContent: {
      maxHeight: 200,
      backgroundColor: colors.background,
      borderRadius: radii.md,
      padding: spacing[2],
    },
  });
}

export function SelectContent({ children, ...props }: SelectContentProps & any) {
  const { isOpen, setIsOpen, value, onValueChange } = props;
  const { colors } = useTheme();
  const styles = useMemo(() => createContentStyles(colors), [colors]);

  return (
    <Modal visible={isOpen} onClose={() => setIsOpen(false)}>
      <View style={styles.selectContent}>
        <ScrollView>
          {React.Children.map(children, (child) =>
            React.isValidElement(child)
              ? React.cloneElement(child as React.ReactElement<any>, {
                  onSelect: (itemValue: string) => {
                    onValueChange(itemValue);
                    setIsOpen(false);
                  },
                  isSelected: value === (child.props as any).value,
                })
              : child
          )}
        </ScrollView>
      </View>
    </Modal>
  );
}

function createItemStyles(colors: ColorPalette) {
  return StyleSheet.create({
    selectItem: {
      paddingVertical: spacing[2],
      paddingHorizontal: spacing[3],
      borderRadius: radii.sm,
    },
    selectItemSelected: { backgroundColor: colors.muted },
    selectItemText: { fontSize: fontSize.base, color: colors.foreground },
    selectItemTextSelected: { fontWeight: fontWeights.medium },
  });
}

export function SelectItem({ value, children, onSelect, isSelected }: SelectItemProps & any) {
  const { colors } = useTheme();
  const styles = useMemo(() => createItemStyles(colors), [colors]);

  return (
    <Pressable
      style={[styles.selectItem, isSelected && styles.selectItemSelected]}
      onPress={() => onSelect(value)}
    >
      <Text style={[styles.selectItemText, isSelected && styles.selectItemTextSelected]}>
        {children}
      </Text>
    </Pressable>
  );
}
