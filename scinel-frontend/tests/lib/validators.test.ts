import { describe, it, expect } from '@jest/globals';
import { isValidEmail, isValidPassword, isStrongPassword } from '@/lib/validators';

describe('validators', () => {
  it('should validate email', () => {
    expect(isValidEmail('test@example.com')).toBe(true);
    expect(isValidEmail('invalid-email')).toBe(false);
    expect(isValidEmail('test@')).toBe(false);
  });

  it('should validate password length', () => {
    expect(isValidPassword('password123')).toBe(true);
    expect(isValidPassword('pass')).toBe(false);
  });

  it('should check strong password', () => {
    expect(isStrongPassword('Password123!')).toBe(true);
    expect(isStrongPassword('password')).toBe(false);
  });
});
