import React, { useState } from 'react';
import { View, Text, Pressable, StyleSheet, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Modal } from './modal';
import { colors, spacing, radii, fontSize } from './tokens';

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
    <View style={styles.select}>
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

export function SelectTrigger({ children, ...props }: SelectTriggerProps & any) {
  const { isOpen, setIsOpen } = props;

  return (
    <Pressable
      style={styles.selectTrigger}
      onPress={() => setIsOpen(!isOpen)}
    >
      {children}
      <Ionicons
        name={isOpen ? 'chevron-up' : 'chevron-down'}
        size={16}
        color={colors['muted-foreground']}
      />
    </Pressable>
  );
}

export function SelectValue({ placeholder }: SelectValueProps) {
  return <Text style={styles.selectValue}>{placeholder}</Text>;
}

export function SelectContent({ children, ...props }: SelectContentProps & any) {
  const { isOpen, setIsOpen, value, onValueChange } = props;

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

export function SelectItem({ value, children, onSelect, isSelected }: SelectItemProps & any) {
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

const styles = StyleSheet.create({
  select: {
    position: 'relative',
  },
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
  selectValue: {
    fontSize: fontSize.base,
    color: colors.foreground,
  },
  selectContent: {
    maxHeight: 200,
    backgroundColor: colors.background,
    borderRadius: radii.md,
    padding: spacing[2],
  },
  selectItem: {
    paddingVertical: spacing[2],
    paddingHorizontal: spacing[3],
    borderRadius: radii.sm,
  },
  selectItemSelected: {
    backgroundColor: colors.muted,
  },
  selectItemText: {
    fontSize: fontSize.base,
    color: colors.foreground,
  },
  selectItemTextSelected: {
    fontWeight: '500',
  },
});