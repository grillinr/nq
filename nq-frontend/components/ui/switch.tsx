import React from 'react';
import { Switch as RNSwitch, SwitchProps as RNSwitchProps } from 'react-native';
import { colors } from './tokens';

function Switch(props: Omit<RNSwitchProps, 'trackColor' | 'thumbColor'>) {
  return (
    <RNSwitch
      trackColor={{ false: colors.muted, true: colors.primary }}
      thumbColor={props.value ? colors.primary : colors['muted-foreground']}
      {...props}
    />
  );
}

export default Switch;