import { flattenStyles } from '../utils';

describe('flattenStyles', () => {
  it('should merge single style object', () => {
    const style = { backgroundColor: 'red', padding: 10 };
    const result = flattenStyles([style]);

    expect(result).toEqual(style);
  });

  it('should merge multiple style objects', () => {
    const style1 = { backgroundColor: 'red', padding: 10 };
    const style2 = { margin: 5, fontSize: 16 };

    const result = flattenStyles([style1, style2]);

    expect(result).toMatchObject({
      backgroundColor: 'red',
      padding: 10,
      margin: 5,
      fontSize: 16,
    });
  });

  it('should override conflicting styles (later wins)', () => {
    const style1 = { backgroundColor: 'red', padding: 10 };
    const style2 = { backgroundColor: 'blue' };

    const result = flattenStyles([style1, style2]);

    expect(result).toMatchObject({
      backgroundColor: 'blue',
      padding: 10,
    });
  });

  it('should filter out undefined values', () => {
    const style1 = { backgroundColor: 'red' };
    const result = flattenStyles([style1, undefined, null]);

    expect(result).toEqual({ backgroundColor: 'red' });
  });

  it('should filter out false values', () => {
    const style1 = { backgroundColor: 'red' };
    const style2 = false;

    const result = flattenStyles([style1, style2]);

    expect(result).toEqual({ backgroundColor: 'red' });
  });

  it('should handle empty array', () => {
    const result = flattenStyles([]);
    expect(result).toEqual({});
  });

  it('should handle all falsy values', () => {
    const result = flattenStyles([undefined, null, false]);
    expect(result).toEqual({});
  });

  it('should handle conditional styles pattern', () => {
    const isActive = true;
    const isDisabled = false;

    const result = flattenStyles([
      { backgroundColor: 'white' },
      isActive && { borderColor: 'blue', borderWidth: 2 },
      isDisabled && { opacity: 0.5 },
    ]);

    expect(result).toMatchObject({
      backgroundColor: 'white',
      borderColor: 'blue',
      borderWidth: 2,
    });
    expect(result).not.toHaveProperty('opacity');
  });
});
