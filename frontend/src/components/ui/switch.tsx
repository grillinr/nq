import React from 'react';
import { Switch as RNSwitch, SwitchProps as RNSwitchProps } from 'react-native';
import { useTheme } from './theme-provider';

function Switch({ value, ...rest }: Omit<RNSwitchProps, 'trackColor' | 'thumbColor'>) {
  const { colors } = useTheme();
  return (
    <RNSwitch
      trackColor={{ false: colors.muted, true: colors.secondary }}
      thumbColor={value ? colors.background : colors.mutedForeground}
      value={value}
      {...rest}
    />
  );
}

export default Switch;
