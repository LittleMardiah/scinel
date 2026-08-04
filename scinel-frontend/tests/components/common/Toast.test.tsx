import Toast from '@/components/common/Toast';
import { describe, it, expect, jest, beforeEach } from '@jest/globals';

describe('Toast', () => {
  beforeEach(() => {
    // Reset messages
    (Toast as any).getMessages = () => [];
    // Clear listeners
    (Toast as any).listeners = [];
  });

  it('should add toast message', () => {
    const message = 'Test message';
    Toast.success(message);
    const messages = Toast.getMessages();
    expect(messages.length).toBe(1);
    expect(messages[0].message).toBe(message);
    expect(messages[0].type).toBe('success');
  });

  it('should add error toast', () => {
    Toast.error('Error message');
    const messages = Toast.getMessages();
    expect(messages[0].type).toBe('error');
  });

  it('should add warning toast', () => {
    Toast.warning('Warning message');
    const messages = Toast.getMessages();
    expect(messages[0].type).toBe('warning');
  });

  it('should add info toast', () => {
    Toast.info('Info message');
    const messages = Toast.getMessages();
    expect(messages[0].type).toBe('info');
  });
});
