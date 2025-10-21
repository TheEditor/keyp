/**
 * Unit tests for crypto module
 */

import { test } from 'node:test';
import * as assert from 'node:assert/strict';
import { encrypt, decrypt, deriveKey, verifyPassword, EncryptionResult } from './crypto.js';

test('Crypto Module', async (t) => {
  await t.test('deriveKey generates consistent keys for same password and salt', () => {
    const password = 'test-password';
    const { key: key1, salt } = deriveKey(password);
    const { key: key2 } = deriveKey(password, salt, 100000);

    assert.deepEqual(key1, key2, 'Keys should match when using same salt and iterations');
  });

  await t.test('deriveKey generates different keys for different passwords', () => {
    const { key: key1, salt } = deriveKey('password1');
    const { key: key2 } = deriveKey('password2', salt, 100000);

    assert.notDeepEqual(key1, key2, 'Keys should differ for different passwords');
  });

  await t.test('deriveKey generates different salt on each call', () => {
    const { salt: salt1 } = deriveKey('password');
    const { salt: salt2 } = deriveKey('password');

    assert.notDeepEqual(salt1, salt2, 'Salts should be different');
  });

  await t.test('deriveKey requires minimum iterations for security', () => {
    assert.throws(
      () => deriveKey('password', undefined, 50000),
      /PBKDF2 iterations must be at least 100,000/,
      'Should reject iterations < 100,000'
    );
  });

  await t.test('encrypt creates valid encryption result', () => {
    const result = encrypt('secret data', 'password');

    assert.ok(result.ciphertext, 'Should have ciphertext');
    assert.ok(result.authTag, 'Should have authTag');
    assert.ok(result.iv, 'Should have iv');
    assert.ok(result.salt, 'Should have salt');

    // Check that they're base64 encoded
    assert.ok(Buffer.from(result.ciphertext, 'base64').length > 0);
    assert.ok(Buffer.from(result.authTag, 'base64').length > 0);
    assert.ok(Buffer.from(result.iv, 'base64').length > 0);
    assert.ok(Buffer.from(result.salt, 'base64').length > 0);
  });

  await t.test('encrypt produces different ciphertext for same plaintext', () => {
    const result1 = encrypt('secret data', 'password');
    const result2 = encrypt('secret data', 'password');

    assert.notEqual(result1.ciphertext, result2.ciphertext, 'Should produce different ciphertexts (due to random IV and salt)');
  });

  await t.test('decrypt recovers original plaintext', () => {
    const plaintext = 'secret data 123!@#';
    const password = 'my-password';

    const encrypted = encrypt(plaintext, password);
    const decrypted = decrypt(encrypted, password);

    assert.equal(decrypted, plaintext, 'Decrypted text should match original');
  });

  await t.test('decrypt fails with wrong password', () => {
    const encrypted = encrypt('secret data', 'correct-password');

    assert.throws(
      () => decrypt(encrypted, 'wrong-password'),
      /Decryption failed/,
      'Should fail with wrong password'
    );
  });

  await t.test('decrypt handles empty string', () => {
    const plaintext = '';
    const encrypted = encrypt(plaintext, 'password');
    const decrypted = decrypt(encrypted, 'password');

    assert.equal(decrypted, plaintext, 'Should handle empty strings');
  });

  await t.test('decrypt handles large data', () => {
    const largeData = 'x'.repeat(100000);
    const encrypted = encrypt(largeData, 'password');
    const decrypted = decrypt(encrypted, 'password');

    assert.equal(decrypted, largeData, 'Should handle large data');
  });

  await t.test('decrypt handles special characters', () => {
    const plaintext = '{"key": "value", "emoji": "🔒", "unicode": "你好"}';
    const encrypted = encrypt(plaintext, 'password');
    const decrypted = decrypt(encrypted, 'password');

    assert.equal(decrypted, plaintext, 'Should preserve special characters');
  });

  await t.test('verifyPassword returns true for correct password', () => {
    const encrypted = encrypt('data', 'correct-password');
    const isCorrect = verifyPassword(encrypted, 'correct-password');

    assert.equal(isCorrect, true, 'Should return true for correct password');
  });

  await t.test('verifyPassword returns false for incorrect password', () => {
    const encrypted = encrypt('data', 'correct-password');
    const isCorrect = verifyPassword(encrypted, 'wrong-password');

    assert.equal(isCorrect, false, 'Should return false for wrong password');
  });

  await t.test('encrypt with custom iterations', () => {
    const plaintext = 'secret';
    const encrypted = encrypt(plaintext, 'password', 150000);
    const decrypted = decrypt(encrypted, 'password', 150000);

    assert.equal(decrypted, plaintext, 'Should work with custom iterations');
  });

  await t.test('decrypt fails with corrupted ciphertext', () => {
    const encrypted = encrypt('secret', 'password');
    const corrupted: EncryptionResult = {
      ...encrypted,
      ciphertext: Buffer.from('corrupted-data').toString('base64'),
    };

    assert.throws(
      () => decrypt(corrupted, 'password'),
      /Decryption failed/,
      'Should fail with corrupted ciphertext'
    );
  });

  await t.test('decrypt fails with corrupted auth tag', () => {
    const encrypted = encrypt('secret', 'password');
    const corrupted: EncryptionResult = {
      ...encrypted,
      authTag: Buffer.from('bad-tag').toString('base64'),
    };

    assert.throws(
      () => decrypt(corrupted, 'password'),
      /Decryption failed/,
      'Should fail with corrupted auth tag'
    );
  });

  await t.test('encrypt and decrypt preserve JSON structure', () => {
    const jsonData = {
      name: 'test',
      secrets: {
        api_key: 'sk-123',
        token: 'tk-456',
      },
      count: 42,
    };
    const plaintext = JSON.stringify(jsonData);

    const encrypted = encrypt(plaintext, 'password');
    const decrypted = decrypt(encrypted, 'password');
    const parsed = JSON.parse(decrypted);

    assert.deepEqual(parsed, jsonData, 'Should preserve JSON structure');
  });
});
