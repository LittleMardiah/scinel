# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: scinel-e2e.spec.ts >> ScIntel E2E Flows >> Flow 5: Toggle Dark Mode → Persists
- Location: tests/e2e/scinel-e2e.spec.ts:133:7

# Error details

```
Test timeout of 30000ms exceeded.
```

```
Error: page.fill: Test timeout of 30000ms exceeded.
Call log:
  - waiting for locator('input[name="email"]')

```

# Page snapshot

```yaml
- generic [active] [ref=e1]: Cannot GET /login
```

# Test source

```ts
  36  |     await page.fill('input[name="email"]', 'test@scinel.com');
  37  |     await page.fill('input[name="password"]', 'SecurePass123!');
  38  |     await page.click('button[type="submit"]');
  39  |     await page.waitForURL(`${BASE_URL}/dashboard`);
  40  | 
  41  |     // Navigate to new transaction
  42  |     await page.goto(`${BASE_URL}/inventory/transaction/new`);
  43  | 
  44  |     // Fill form
  45  |     await page.fill('input[name="sku_code"]', 'SKU-NIKE-001');
  46  |     await page.selectOption('select[name="type"]', 'in');
  47  |     await page.fill('input[name="quantity"]', '5');
  48  |     await page.selectOption('select[name="warehouse_id"]', '62c35799-3209-4f64-85fc-4103557f5295');
  49  |     await page.fill('input[name="reason"]', 'Automated E2E test');
  50  | 
  51  |     // Submit
  52  |     await page.click('button[type="submit"]');
  53  | 
  54  |     // Verify success toast
  55  |     await expect(page.locator('[role="status"]')).toContainText('Transaksi berhasil dibuat');
  56  | 
  57  |     // Navigate to transaction history
  58  |     await page.goto(`${BASE_URL}/inventory/transactions`);
  59  | 
  60  |     // Verify transaction appears (at least one row)
  61  |     const rows = page.locator('table tbody tr');
  62  |     await expect(rows.first()).toBeVisible();
  63  |   });
  64  | 
  65  |   test('Flow 3: View Forecast → Check Accuracy', async ({ page }) => {
  66  |     // Login
  67  |     await page.goto(`${BASE_URL}/login`);
  68  |     await page.fill('input[name="email"]', 'test@scinel.com');
  69  |     await page.fill('input[name="password"]', 'SecurePass123!');
  70  |     await page.click('button[type="submit"]');
  71  |     await page.waitForURL(`${BASE_URL}/dashboard`);
  72  | 
  73  |     // Go to forecast
  74  |     await page.goto(`${BASE_URL}/forecast`);
  75  | 
  76  |     // Check chart or at least table
  77  |     const chart = page.locator('canvas, .recharts-wrapper');
  78  |     await expect(chart.first()).toBeVisible({ timeout: 10000 });
  79  | 
  80  |     // Go to accuracy
  81  |     await page.click('a[href="/forecast/accuracy"]');
  82  |     await page.waitForURL(`${BASE_URL}/forecast/accuracy`);
  83  | 
  84  |     // Check accuracy metrics
  85  |     await expect(page.locator('table')).toBeVisible();
  86  |   });
  87  | 
  88  |   test('Flow 4: Acknowledge Anomaly → Update Status', async ({ page }) => {
  89  |     // Login
  90  |     await page.goto(`${BASE_URL}/login`);
  91  |     await page.fill('input[name="email"]', 'test@scinel.com');
  92  |     await page.fill('input[name="password"]', 'SecurePass123!');
  93  |     await page.click('button[type="submit"]');
  94  |     await page.waitForURL(`${BASE_URL}/dashboard`);
  95  | 
  96  |     // Go to anomalies
  97  |     await page.goto(`${BASE_URL}/anomalies`);
  98  | 
  99  |     // Check if there's at least one anomaly row
  100 |     const rows = page.locator('table tbody tr');
  101 |     const count = await rows.count();
  102 | 
  103 |     if (count === 0) {
  104 |       // Skip if no anomalies (create one via API if needed)
  105 |       test.skip(true, 'No anomalies available to test');
  106 |       return;
  107 |     }
  108 | 
  109 |     // Click first anomaly row (to go to detail)
  110 |     await rows.first().click();
  111 |     await page.waitForURL(/\/anomalies\/[a-f0-9-]+/);
  112 | 
  113 |     // Click Acknowledge button
  114 |     const ackBtn = page.locator('button:has-text("Acknowledge")');
  115 |     if (await ackBtn.isVisible()) {
  116 |       await ackBtn.click();
  117 | 
  118 |       // Fill notes if modal appears
  119 |       const notesInput = page.locator('textarea[name="notes"]');
  120 |       if (await notesInput.isVisible()) {
  121 |         await notesInput.fill('Acknowledged via automation');
  122 |         await page.click('button:has-text("Confirm")');
  123 |       }
  124 | 
  125 |       // Verify status changed
  126 |       await expect(page.locator('[data-testid="status"]')).toContainText('acknowledged');
  127 |     } else {
  128 |       // Maybe already resolved, skip
  129 |       test.skip(true, 'Anomaly already acknowledged or resolved');
  130 |     }
  131 |   });
  132 | 
  133 |   test('Flow 5: Toggle Dark Mode → Persists', async ({ page }) => {
  134 |     // Login
  135 |     await page.goto(`${BASE_URL}/login`);
> 136 |     await page.fill('input[name="email"]', 'test@scinel.com');
      |                ^ Error: page.fill: Test timeout of 30000ms exceeded.
  137 |     await page.fill('input[name="password"]', 'SecurePass123!');
  138 |     await page.click('button[type="submit"]');
  139 |     await page.waitForURL(`${BASE_URL}/dashboard`);
  140 | 
  141 |     // Get initial theme (should be light by default)
  142 |     const initialHtmlClass = await page.locator('html').getAttribute('class');
  143 |     const wasDark = initialHtmlClass?.includes('dark') || false;
  144 | 
  145 |     // Click theme toggle
  146 |     await page.click('[data-testid="theme-toggle"]');
  147 | 
  148 |     // Verify theme changed
  149 |     const newHtmlClass = await page.locator('html').getAttribute('class');
  150 |     const isDark = newHtmlClass?.includes('dark') || false;
  151 | 
  152 |     // Theme should be toggled
  153 |     expect(isDark).toBe(!wasDark);
  154 | 
  155 |     // Refresh page
  156 |     await page.reload();
  157 | 
  158 |     // Verify theme persists
  159 |     const afterRefreshClass = await page.locator('html').getAttribute('class');
  160 |     const afterRefreshIsDark = afterRefreshClass?.includes('dark') || false;
  161 |     expect(afterRefreshIsDark).toBe(isDark);
  162 |   });
  163 | });
```