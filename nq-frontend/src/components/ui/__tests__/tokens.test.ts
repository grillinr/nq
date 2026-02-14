import { lightColors, darkColors, sharedColors, fontWeights, radii, spacing, fontSize } from '../tokens';

describe('Design Tokens', () => {
  describe('lightColors', () => {
    it('should have all required color properties', () => {
      expect(lightColors).toHaveProperty('background');
      expect(lightColors).toHaveProperty('foreground');
      expect(lightColors).toHaveProperty('primary');
      expect(lightColors).toHaveProperty('secondary');
    });

    it('should have valid hex color values for primary colors', () => {
      expect(lightColors.background).toMatch(/^#[0-9a-f]{6}$/i);
      expect(lightColors.foreground).toMatch(/^#[0-9a-f]{6}$/i);
      expect(lightColors.primary).toMatch(/^#[0-9a-f]{6}$/i);
    });

    it('should have accessible contrast between background and foreground', () => {
      expect(lightColors.background).not.toBe(lightColors.foreground);
    });
  });

  describe('darkColors', () => {
    it('should have all required color properties', () => {
      expect(darkColors).toHaveProperty('background');
      expect(darkColors).toHaveProperty('foreground');
      expect(darkColors).toHaveProperty('primary');
      expect(darkColors).toHaveProperty('secondary');
    });

    it('should have valid color values for primary colors', () => {
      expect(darkColors.background).toBeTruthy();
      expect(darkColors.foreground).toBeTruthy();
      expect(darkColors.primary).toBeTruthy();
    });

    it('should differ from light theme colors', () => {
      expect(darkColors.background).not.toBe(lightColors.background);
      expect(darkColors.foreground).not.toBe(lightColors.foreground);
    });
  });

  describe('sharedColors', () => {
    it('should have shared color properties', () => {
      expect(sharedColors).toHaveProperty('primary');
      expect(sharedColors).toHaveProperty('primaryForeground');
      expect(sharedColors).toHaveProperty('input');
      expect(sharedColors).toHaveProperty('chart3');
      expect(sharedColors).toHaveProperty('chart4');
    });

    it('should match primary in both light and dark themes', () => {
      expect(lightColors.primary).toBe(sharedColors.primary);
      expect(darkColors.primary).toBe(sharedColors.primary);
      expect(lightColors.primaryForeground).toBe(sharedColors.primaryForeground);
      expect(darkColors.primaryForeground).toBe(sharedColors.primaryForeground);
    });
  });

  describe('fontWeights', () => {
    it('should have medium and normal weights', () => {
      expect(fontWeights.medium).toBe('500');
      expect(fontWeights.normal).toBe('400');
    });

    it('should have string values for font weights', () => {
      expect(typeof fontWeights.medium).toBe('string');
      expect(typeof fontWeights.normal).toBe('string');
    });
  });

  describe('radii', () => {
    it('should have all size variants', () => {
      expect(radii).toHaveProperty('sm');
      expect(radii).toHaveProperty('md');
      expect(radii).toHaveProperty('lg');
      expect(radii).toHaveProperty('xl');
    });

    it('should have numeric values', () => {
      expect(typeof radii.sm).toBe('number');
      expect(typeof radii.md).toBe('number');
      expect(typeof radii.lg).toBe('number');
      expect(typeof radii.xl).toBe('number');
    });

    it('should have ascending values', () => {
      expect(radii.sm).toBeLessThan(radii.md);
      expect(radii.md).toBeLessThan(radii.lg);
      expect(radii.lg).toBeLessThan(radii.xl);
    });
  });

  describe('spacing', () => {
    it('should have all spacing scale values', () => {
      expect(spacing).toHaveProperty('1');
      expect(spacing).toHaveProperty('2');
      expect(spacing).toHaveProperty('3');
      expect(spacing).toHaveProperty('4');
      expect(spacing).toHaveProperty('6');
      expect(spacing).toHaveProperty('8');
    });

    it('should have numeric values', () => {
      expect(typeof spacing[1]).toBe('number');
      expect(typeof spacing[2]).toBe('number');
    });

    it('should have consistent spacing scale', () => {
      expect(spacing[2]).toBe(spacing[1] * 2);
      expect(spacing[4]).toBe(spacing[2] * 2);
      expect(spacing[8]).toBe(spacing[4] * 2);
    });
  });

  describe('fontSize', () => {
    it('should have all font size variants', () => {
      expect(fontSize).toHaveProperty('xs');
      expect(fontSize).toHaveProperty('sm');
      expect(fontSize).toHaveProperty('base');
      expect(fontSize).toHaveProperty('lg');
      expect(fontSize).toHaveProperty('xl');
      expect(fontSize).toHaveProperty('2xl');
    });

    it('should have numeric values', () => {
      expect(typeof fontSize.xs).toBe('number');
      expect(typeof fontSize.sm).toBe('number');
      expect(typeof fontSize.base).toBe('number');
    });

    it('should have ascending values', () => {
      expect(fontSize.xs).toBeLessThan(fontSize.sm);
      expect(fontSize.sm).toBeLessThan(fontSize.base);
      expect(fontSize.base).toBeLessThan(fontSize.lg);
      expect(fontSize.lg).toBeLessThan(fontSize.xl);
      expect(fontSize.xl).toBeLessThan(fontSize['2xl']);
    });
  });
});
