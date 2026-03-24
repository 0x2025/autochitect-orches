const { test } = require('@jest/globals');

test('basic test', () => {
  const assert = require('assert');
  assert.strictEqual(1, 1, '1 should equal 1');
});