import { logError, logInfo } from '../logger';

describe('logger', () => {
  let consoleSpy: jest.SpyInstance;
  let consoleErrorSpy: jest.SpyInstance;
  const originalEnv = process.env.NODE_ENV;

  beforeEach(() => {
    consoleSpy = jest.spyOn(console, 'log').mockImplementation();
    consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
  });

  afterEach(() => {
    consoleSpy.mockRestore();
    consoleErrorSpy.mockRestore();
    process.env.NODE_ENV = originalEnv;
  });

  describe('logError', () => {
    it('should log errors in development', () => {
      process.env.NODE_ENV = 'development';
      logError('test error', { foo: 'bar' });
      
      expect(consoleErrorSpy).toHaveBeenCalledWith('test error', { foo: 'bar' });
    });

    it('should not log errors in production', () => {
      process.env.NODE_ENV = 'production';
      logError('test error');
      
      expect(consoleErrorSpy).not.toHaveBeenCalled();
    });

    it('should handle multiple arguments', () => {
      process.env.NODE_ENV = 'test';
      const args = ['error1', 'error2', { code: 500 }];
      logError(...args);
      
      expect(consoleErrorSpy).toHaveBeenCalledWith(...args);
    });
  });

  describe('logInfo', () => {
    it('should log info in development', () => {
      process.env.NODE_ENV = 'development';
      logInfo('test info', { data: 'value' });
      
      expect(consoleSpy).toHaveBeenCalledWith('test info', { data: 'value' });
    });

    it('should not log info in production', () => {
      process.env.NODE_ENV = 'production';
      logInfo('test info');
      
      expect(consoleSpy).not.toHaveBeenCalled();
    });

    it('should log info in test environment', () => {
      process.env.NODE_ENV = 'test';
      logInfo('test message');
      
      expect(consoleSpy).toHaveBeenCalledWith('test message');
    });

    it('should handle multiple arguments', () => {
      process.env.NODE_ENV = 'development';
      const args = ['info1', 'info2', { status: 'ok' }];
      logInfo(...args);
      
      expect(consoleSpy).toHaveBeenCalledWith(...args);
    });
  });
});
