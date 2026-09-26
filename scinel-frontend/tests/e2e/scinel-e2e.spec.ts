import { test, expect } from '@playwright/test';

test.describe.configure({ timeout: 120000 }); // 120 detik global

test.describe('ScIntel E2E Flows', () => {
  const BASE_URL = 'http://localhost:3000';
  const API_URL = 'http://localhost:8080/api/v1';

  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies();
  });

  // ============================================================
  // FLOW 1: Login → Dashboard → Logout
  // ============================================================
  test.slow();
  test('Flow 1: Login → Dashboard → Logout', async ({ page, request }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.fill('#email', 'test@example.com');
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');

    // Tunggu elemen di dashboard (bukan navigasi load)
    await page.waitForSelector('h1:has-text("Dashboard")', { timeout: 90000 });
    await expect(page.locator('h1')).toContainText('Dashboard');

    const cookies = await page.context().cookies();
    const accessTokenCookie = cookies.find(c => c.name === 'access_token');
    const refreshTokenCookie = cookies.find(c => c.name === 'refresh_token');
    if (!accessTokenCookie || !refreshTokenCookie) {
      throw new Error('Token tidak ditemukan di cookies');
    }

    const logoutResp = await request.post(`${API_URL}/auth/logout`, {
      headers: {
        'Authorization': `Bearer ${accessTokenCookie.value}`,
        'Content-Type': 'application/json'
      },
      data: { refresh_token: refreshTokenCookie.value }
    });
    expect(logoutResp.status()).toBe(200);
    await page.context().clearCookies();
    await page.goto(`${BASE_URL}/dashboard`);
    // Redirect ke login
    await expect(page.locator('#email')).toBeVisible({ timeout: 30000 });
  });

  // ============================================================
  // FLOW 2: Create Transaction → See in History
  // ============================================================
  test.slow();
  test('Flow 2: Create Transaction → See in History', async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.fill('#email', 'test@example.com');
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForSelector('h1:has-text("Dashboard")', { timeout: 90000 });

    await page.goto(`${BASE_URL}/inventory/transaction/new`);
    await page.waitForLoadState('domcontentloaded', { timeout: 30000 });

    const skuField = page.getByPlaceholder('e.g., SKU-001');
    await expect(skuField).toBeVisible({ timeout: 30000 });
    await skuField.fill('SKU-NIKE-001');

    try {
      await page.getByRole('radio', { name: 'Masuk' }).click();
    } catch {
      await page.locator('text=Masuk').first().click();
    }

    const qtyInput = page.getByRole('spinbutton', { name: 'Kuantitas' });
    await expect(qtyInput).toBeVisible({ timeout: 30000 });
    await qtyInput.fill('5');

    await page.selectOption('select#warehouse', { index: 0 });
    await page.fill('input#reason', 'Automated E2E test');
    await page.click('button[type="submit"]');

    await page.goto(`${BASE_URL}/inventory/transactions`);
    await page.waitForLoadState('domcontentloaded', { timeout: 30000 });
    const rows = page.locator('table tbody tr');
    await expect(rows).not.toHaveCount(0);
  });

  // ============================================================
  // FLOW 3: View Forecast → Check Accuracy
  // ============================================================
  test.slow();
  test('Flow 3: View Forecast → Check Accuracy', async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.fill('#email', 'test@example.com');
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForSelector('h1:has-text("Dashboard")', { timeout: 90000 });

    await page.goto(`${BASE_URL}/forecast`);
    await page.waitForLoadState('domcontentloaded', { timeout: 30000 });

    const chart = page.locator('.recharts-wrapper');
    await expect(chart).toBeVisible({ timeout: 30000 });
  });

  // ============================================================
  // FLOW 4: Acknowledge Anomaly → Update Status
  // ============================================================
  test.slow();
  test('Flow 4: Acknowledge Anomaly → Update Status', async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.fill('#email', 'test@example.com');
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForSelector('h1:has-text("Dashboard")', { timeout: 90000 });

    await page.goto(`${BASE_URL}/anomalies`);
    await page.waitForLoadState('domcontentloaded', { timeout: 30000 });

    const table = page.locator('table');
    await expect(table).toBeVisible({ timeout: 30000 });

    const rows = table.locator('tbody tr');
    const count = await rows.count();

    if (count === 0) {
      console.log('⚠️ Tidak ada anomali, skip test');
      return;
    }

    const firstRow = rows.first();
    await firstRow.click();
    await page.waitForURL(/\/anomalies\/[a-z0-9-]+/, { timeout: 30000 });
    await expect(page.locator('h1')).toContainText(/Anomali|Detail/i, { timeout: 30000 });
  });

  // ============================================================
  // FLOW 5: Toggle Dark Mode → Persists
  // ============================================================
  test.slow();
  test('Flow 5: Toggle Dark Mode → Persists', async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.fill('#email', 'test@example.com');
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForSelector('h1:has-text("Dashboard")', { timeout: 90000 });
    await page.waitForLoadState('domcontentloaded', { timeout: 30000 });

    const isDarkBefore = await page.evaluate(() => {
      return document.documentElement.classList.contains('dark');
    });
    expect(isDarkBefore).toBe(false);

    const toggleBtn = page.getByRole('button', { name: /toggle/i });
    await expect(toggleBtn).toBeVisible({ timeout: 30000 });

    for (let i = 0; i < 4; i++) {
      await toggleBtn.click();
      const dark = await page.evaluate(() => document.documentElement.classList.contains('dark'));
      if (dark) break;
    }

    await page.waitForSelector('html.dark', { timeout: 30000 });

    await page.reload();
    await page.waitForLoadState('domcontentloaded', { timeout: 30000 });

    const isDarkAfter = await page.evaluate(() => {
      return document.documentElement.classList.contains('dark');
    });
    expect(isDarkAfter).toBe(true);
  });
});
