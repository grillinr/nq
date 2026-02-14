import { useCallback, useState, useRef } from "react";
import { NativeSyntheticEvent, NativeScrollEvent } from "react-native";

interface UseScrollHeaderResult {
  isHeaderVisible: boolean;
  handleScroll: (event: NativeSyntheticEvent<NativeScrollEvent>) => void;
}

export function useScrollHeader(threshold: number = 50): UseScrollHeaderResult {
  const [isHeaderVisible, setIsHeaderVisible] = useState(true);
  const lastScrollY = useRef(0);

  const handleScroll = useCallback(
    (event: NativeSyntheticEvent<NativeScrollEvent>) => {
      const currentScrollY = event.nativeEvent.contentOffset.y;

      // Don't hide if we're near the top
      if (currentScrollY < threshold) {
        setIsHeaderVisible(true);
      } else {
        // Hide when scrolling down, show when scrolling up
        if (currentScrollY > lastScrollY.current) {
          // Scrolling down
          setIsHeaderVisible(false);
        } else if (currentScrollY < lastScrollY.current) {
          // Scrolling up
          setIsHeaderVisible(true);
        }
      }

      lastScrollY.current = currentScrollY;
    },
    [threshold]
  );

  return { isHeaderVisible, handleScroll };
}
