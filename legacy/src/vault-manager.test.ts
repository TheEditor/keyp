/**
 * Integration tests for vault manager
 */

import { test } from 'node:test';
import * as assert from 'node:assert/strict';
import { rmSync, existsSync, mkdirSync } from 'fs';
import { join } from 'path';
import { tmpdir } from 'os';
import { VaultManager } from './vault-manager.js';
import { SecretsManager } from './secrets.js';

// Create a temporary directory for test vaults
const testVaultDir = join(tmpdir(), `keyp-test-${Date.now()}`);
mkdirSync(testVaultDir, { recursive: true });

test('VaultManager Integration Tests', async (t) => {
  const getTestVaultPath = () => {
    const vaultPath = join(testVaultDir, `vault-${Math.random()}.json`);
    return vaultPath;
  };

  t.after(() => {
    // Cleanup test vaults
    if (existsSync(testVaultDir)) {
      rmSync(testVaultDir, { recursive: true, force: true });
    }
  });

  await t.test('Initialize new vault', async () => {
    const vaultPath = getTestVaultPath();
    const manager = new VaultManager(vaultPath);

    const result = await manager.initializeVault('test-password');

    assert.equal(result.success, true, 'Should initialize successfully');
    assert.equal(manager.vaultFileExists(), true, 'Vault file should exist');
    assert.equal(manager.isVaultUnlocked(), true, 'Vault should be unlocked after init');
  });

  await t.test('Initialize fails if vault already exists', async () => {
    const vaultPath = getTestVaultPath();
    const manager = new VaultManager(vaultPath);

    await manager.initializeVault('password1');
    const result2 = await manager.initializeVault('password2');

    assert.equal(result2.success, false, 'Should fail to re-initialize');
    assert.match(result2.error || '', /already exists/, 'Error should mention existing vault');
  });

  await t.test('Unlock existing vault with correct password', async () => {
    const vaultPath = getTestVaultPath();
    const manager1 = new VaultManager(vaultPath);
    await manager1.initializeVault('test-password');
    manager1.lockVault();

    const manager2 = new VaultManager(vaultPath);
    const result = await manager2.unlockVault('test-password');

    assert.equal(result.success, true, 'Should unlock successfully');
    assert.equal(manager2.isVaultUnlocked(), true, 'Should be unlocked');
  });

  await t.test('Unlock fails with incorrect password', async () => {
    const vaultPath = getTestVaultPath();
    const manager1 = new VaultManager(vaultPath);
    await manager1.initializeVault('correct-password');
    manager1.lockVault();

    const manager2 = new VaultManager(vaultPath);
    const result = await manager2.unlockVault('wrong-password');

    assert.equal(result.success, false, 'Should fail to unlock');
    assert.equal(manager2.isVaultUnlocked(), false, 'Should remain locked');
  });

  await t.test('Unlock fails if vault does not exist', async () => {
    const vaultPath = getTestVaultPath();
    const manager = new VaultManager(vaultPath);
    const result = await manager.unlockVault('password');

    assert.equal(result.success, false, 'Should fail');
    assert.match(result.error || '', /does not exist/, 'Error should mention non-existent vault');
  });

  await t.test('Lock clears in-memory data', async () => {
    const vaultPath = getTestVaultPath();
    const manager = new VaultManager(vaultPath);
    await manager.initializeVault('password');

    assert.ok(manager.getUnlockedData(), 'Should have data when unlocked');

    manager.lockVault();

    assert.equal(manager.getUnlockedData(), null, 'Should have no data when locked');
    assert.equal(manager.isVaultUnlocked(), false, 'Should be locked');
  });

  await t.test('Save and reload vault', async () => {
    const vaultPath = getTestVaultPath();
    const password = 'test-password';

    // Create and save
    const manager1 = new VaultManager(vaultPath);
    await manager1.initializeVault(password);
    const data1 = manager1.getUnlockedData();
    if (data1) {
      data1['api-key'] = 'sk-test-123';
      data1['token'] = 'tk-test-456';
    }
    await manager1.saveVault(password);
    manager1.lockVault();

    // Reload and verify
    const manager2 = new VaultManager(vaultPath);
    await manager2.unlockVault(password);
    const data2 = manager2.getUnlockedData();

    assert.deepEqual(data2, { 'api-key': 'sk-test-123', token: 'tk-test-456' }, 'Should restore saved secrets');
  });

  await t.test('Save fails when vault is locked', async () => {
    const vaultPath = getTestVaultPath();
    const manager = new VaultManager(vaultPath);
    await manager.initializeVault('password');
    manager.lockVault();

    const result = await manager.saveVault('password');

    assert.equal(result.success, false, 'Should fail to save locked vault');
    assert.match(result.error || '', /not unlocked/, 'Error should indicate vault is locked');
  });

  await t.test('SecretsManager set and get operations', async () => {
    const vaultPath = getTestVaultPath();
    const manager = new VaultManager(vaultPath);
    await manager.initializeVault('password');
    const data = manager.getUnlockedData();

    assert.ok(data, 'Should have vault data');

    SecretsManager.setSecret(data, 'github-token', 'ghp_test123');
    const value = SecretsManager.getSecret(data, 'github-token');

    assert.equal(value, 'ghp_test123', 'Should retrieve set secret');
  });

  await t.test('SecretsManager rejects empty name', () => {
    const data = {};
    const result = SecretsManager.setSecret(data, '', 'value');

    assert.equal(result.success, false, 'Should reject empty name');
    assert.match(result.error || '', /cannot be empty/, 'Error should mention empty name');
  });

  await t.test('SecretsManager rejects empty value', () => {
    const data = {};
    const result = SecretsManager.setSecret(data, 'key', '');

    assert.equal(result.success, false, 'Should reject empty value');
    assert.match(result.error || '', /cannot be empty/, 'Error should mention empty value');
  });

  await t.test('SecretsManager delete operations', () => {
    const data: { [key: string]: string } = {};
    SecretsManager.setSecret(data, 'secret1', 'value1');
    SecretsManager.setSecret(data, 'secret2', 'value2');

    const result = SecretsManager.deleteSecret(data, 'secret1');

    assert.equal(result.success, true, 'Should delete successfully');
    assert.equal(SecretsManager.hasSecret(data, 'secret1'), false, 'Should not have deleted secret');
    assert.equal(SecretsManager.hasSecret(data, 'secret2'), true, 'Should keep other secrets');
  });

  await t.test('SecretsManager delete fails for non-existent secret', () => {
    const data = {};
    const result = SecretsManager.deleteSecret(data, 'non-existent');

    assert.equal(result.success, false, 'Should fail');
    assert.match(result.error || '', /not found/, 'Error should mention not found');
  });

  await t.test('SecretsManager list operations', () => {
    const data: { [key: string]: string } = {};
    SecretsManager.setSecret(data, 'zebra', 'z');
    SecretsManager.setSecret(data, 'apple', 'a');
    SecretsManager.setSecret(data, 'banana', 'b');

    const list = SecretsManager.listSecrets(data);

    assert.deepEqual(list, ['apple', 'banana', 'zebra'], 'Should return sorted list');
  });

  await t.test('SecretsManager getSecretCount', () => {
    const data: { [key: string]: string } = {};
    assert.equal(SecretsManager.getSecretCount(data), 0, 'Should start with 0');

    SecretsManager.setSecret(data, 'secret1', 'value1');
    assert.equal(SecretsManager.getSecretCount(data), 1, 'Should count 1');

    SecretsManager.setSecret(data, 'secret2', 'value2');
    assert.equal(SecretsManager.getSecretCount(data), 2, 'Should count 2');
  });

  await t.test('SecretsManager search operations', () => {
    const data: { [key: string]: string } = {};
    SecretsManager.setSecret(data, 'github-token', 'ghp-123');
    SecretsManager.setSecret(data, 'gitlab-token', 'glpat-456');
    SecretsManager.setSecret(data, 'github-api-key', 'ghp-789');
    SecretsManager.setSecret(data, 'aws-secret', 'aws-123');

    const results = SecretsManager.searchSecrets(data, 'github');

    assert.deepEqual(results, ['github-api-key', 'github-token'], 'Should find matching secrets');
  });

  await t.test('SecretsManager search case-insensitive', () => {
    const data: { [key: string]: string } = {};
    SecretsManager.setSecret(data, 'API-Key', 'value');

    const results = SecretsManager.searchSecrets(data, 'api');

    assert.deepEqual(results, ['API-Key'], 'Should find with case-insensitive search');
  });

  await t.test('SecretsManager clearAllSecrets requires confirmation', () => {
    const data: { [key: string]: string } = {};
    SecretsManager.setSecret(data, 'secret1', 'value1');

    const result = SecretsManager.clearAllSecrets(data);

    assert.equal(result.success, false, 'Should fail without confirmation');
    assert.equal(SecretsManager.getSecretCount(data), 1, 'Should not clear without confirmation');
  });

  await t.test('SecretsManager clearAllSecrets with confirmation', () => {
    const data: { [key: string]: string } = {};
    SecretsManager.setSecret(data, 'secret1', 'value1');
    SecretsManager.setSecret(data, 'secret2', 'value2');

    const result = SecretsManager.clearAllSecrets(data, 'CONFIRM_DELETE_ALL');

    assert.equal(result.success, true, 'Should succeed with confirmation');
    assert.equal(SecretsManager.getSecretCount(data), 0, 'Should clear all secrets');
  });

  await t.test('Full workflow: init, set secrets, save, reload, retrieve', async () => {
    const vaultPath = getTestVaultPath();
    const password = 'my-secure-password';

    // Step 1: Initialize vault
    const manager1 = new VaultManager(vaultPath);
    await manager1.initializeVault(password);

    // Step 2: Add secrets
    const data1 = manager1.getUnlockedData();
    assert.ok(data1, 'Should have data');
    SecretsManager.setSecret(data1, 'db-password', 'super-secret-123');
    SecretsManager.setSecret(data1, 'api-key', 'sk-test-abc');

    // Step 3: Save vault
    const saveResult = await manager1.saveVault(password);
    assert.equal(saveResult.success, true, 'Should save successfully');

    // Step 4: Lock and reload
    manager1.lockVault();
    const manager2 = new VaultManager(vaultPath);
    const unlockResult = await manager2.unlockVault(password);
    assert.equal(unlockResult.success, true, 'Should unlock successfully');

    // Step 5: Verify secrets
    const data2 = manager2.getUnlockedData();
    assert.ok(data2, 'Should have loaded data');
    assert.equal(SecretsManager.getSecret(data2, 'db-password'), 'super-secret-123', 'Should retrieve db-password');
    assert.equal(SecretsManager.getSecret(data2, 'api-key'), 'sk-test-abc', 'Should retrieve api-key');
  });
});
