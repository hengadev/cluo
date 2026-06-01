/**
 * Design review script for the "Changer d'affaire" (CasePicker) feature.
 * Run with: node scripts/review_case_picker.mjs
 */

import { chromium } from 'playwright';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { mkdirSync } from 'fs';

const __dirname = dirname(fileURLToPath(import.meta.url));
const screenshotsDir = join(__dirname, '../.local/review-screenshots');
mkdirSync(screenshotsDir, { recursive: true });

const BASE_URL = 'http://localhost:5173';

async function shot(page, name) {
    const p = join(screenshotsDir, `${name}.png`);
    await page.screenshot({ path: p, fullPage: false });
    console.log(`Screenshot saved: ${p}`);
    return p;
}

(async () => {
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext({
        viewport: { width: 390, height: 844 },   // iPhone 14 Pro ≈ mobile target
        deviceScaleFactor: 2,
        userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15',
    });

    const page = await context.newPage();

    // Capture all console messages
    const consoleLogs = [];
    page.on('console', msg => consoleLogs.push({ type: msg.type(), text: msg.text() }));
    page.on('pageerror', err => consoleLogs.push({ type: 'pageerror', text: err.message }));

    // ── Step 1: Load home page ───────────────────────────────────────────────
    console.log('\n[Step 1] Loading home page…');
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    // Allow Svelte async data (casesPromise, currentCase) to settle
    await page.waitForTimeout(600);
    await shot(page, '01-home-initial');

    // Check that the "Affaire active" section renders
    const affaireActiveHeader = await page.locator('text=Affaire active').isVisible();
    console.log(`  "Affaire active" header visible: ${affaireActiveHeader}`);

    const changerBtn = page.locator('button:has-text("Changer d\'affaire")');
    const changerBtnVisible = await changerBtn.isVisible();
    console.log(`  "Changer d'affaire" button visible: ${changerBtnVisible}`);

    // Grab the initial active case title from CurrentCase
    const currentCaseText = await page.locator('[data-testid="current-case-title"]').textContent().catch(() => null);
    console.log(`  CurrentCase title (data-testid): ${currentCaseText ?? '(no data-testid, inspecting DOM)'}`);

    // ── Step 2: Open the CasePicker bottom sheet ─────────────────────────────
    console.log('\n[Step 2] Clicking "Changer d\'affaire"…');
    await changerBtn.click();
    await page.waitForTimeout(400);   // allow dialog open animation
    await shot(page, '02-case-picker-open');

    const dialogTitle = await page.locator('text=Choisir une affaire').isVisible();
    console.log(`  Dialog title "Choisir une affaire" visible: ${dialogTitle}`);

    // List all case rows rendered in the sheet
    const caseButtons = page.locator('role=dialog >> button');
    const caseCount = await caseButtons.count();
    console.log(`  Case buttons in dialog: ${caseCount}`);
    for (let i = 0; i < caseCount; i++) {
        const txt = (await caseButtons.nth(i).textContent() ?? '').replace(/\s+/g, ' ').trim();
        console.log(`    [${i}] ${txt}`);
    }

    // Active case should have a checkmark
    const checkIcon = page.locator('role=dialog >> [data-lucide="check"]');
    const checkVisible = await checkIcon.isVisible().catch(() => false);
    console.log(`  Check icon on active case visible: ${checkVisible}`);

    // Active case border highlight
    const activeBorder = page.locator('role=dialog >> button.border-dark-900');
    const activeBorderCount = await activeBorder.count();
    console.log(`  Active-border buttons (border-dark-900): ${activeBorderCount}`);

    // ── Step 3: Select a non-active case ─────────────────────────────────────
    // The mock returns mock-case-1 as active; pick the second item (index 1)
    console.log('\n[Step 3] Selecting a non-active case (index 1)…');

    // Find the first case button that does NOT have border-dark-900
    const allCaseButtons = page.locator('role=dialog >> .flex.flex-col.gap-1').locator('..');
    const allCount = await allCaseButtons.count();
    console.log(`  Total case-row containers: ${allCount}`);

    // Click second case button (index 1) — known non-active in mock
    // Use a more reliable selector: the button rows inside the scrollable list
    const listButtons = page.locator('role=dialog >> div.flex.flex-col.gap-2 >> button');
    const listCount = await listButtons.count();
    console.log(`  List buttons found: ${listCount}`);

    let picked = false;
    for (let i = 0; i < listCount; i++) {
        const cls = await listButtons.nth(i).getAttribute('class') ?? '';
        if (!cls.includes('border-dark-900')) {
            const label = (await listButtons.nth(i).textContent() ?? '').replace(/\s+/g, ' ').trim();
            console.log(`  Clicking non-active case at index ${i}: "${label}"`);
            await listButtons.nth(i).click();
            picked = true;
            break;
        }
    }

    if (!picked) {
        console.log('  WARNING: No non-active case found to click!');
    }

    // Wait for dialog to close and UI to update
    await page.waitForTimeout(500);
    await shot(page, '03-after-case-selected');

    // Dialog should be closed
    const dialogClosed = !(await page.locator('text=Choisir une affaire').isVisible());
    console.log(`  Dialog closed after selection: ${dialogClosed}`);

    // ── Step 4: Desktop viewport check ───────────────────────────────────────
    console.log('\n[Step 4] Retesting at desktop width (1440px)…');
    await page.setViewportSize({ width: 1440, height: 900 });
    await page.waitForTimeout(200);
    await shot(page, '04-home-desktop');

    // Re-open picker at desktop
    const changerBtnDesktop = page.locator('button:has-text("Changer d\'affaire")');
    await changerBtnDesktop.click();
    await page.waitForTimeout(300);
    await shot(page, '05-picker-desktop');
    // Close it
    const closeBtn = page.locator('role=dialog >> button').filter({ hasText: '' }).first();
    await page.keyboard.press('Escape');
    await page.waitForTimeout(200);

    // ── Step 5: Tablet viewport ──────────────────────────────────────────────
    console.log('\n[Step 5] Tablet viewport (768px)…');
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(200);
    await shot(page, '06-home-tablet');

    // ── Step 6: Keyboard accessibility ──────────────────────────────────────
    console.log('\n[Step 6] Keyboard accessibility…');
    await page.setViewportSize({ width: 390, height: 844 });
    await page.reload({ waitUntil: 'networkidle' });
    await page.waitForTimeout(400);

    // Tab to the "Changer d'affaire" button
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.waitForTimeout(100);
    await shot(page, '07-keyboard-focus');

    // Focus the button directly and activate with Enter
    await changerBtn.focus();
    await page.keyboard.press('Enter');
    await page.waitForTimeout(400);
    const dialogOpenViaKeyboard = await page.locator('text=Choisir une affaire').isVisible();
    console.log(`  Dialog opens via Enter key: ${dialogOpenViaKeyboard}`);
    await shot(page, '08-dialog-keyboard-open');

    // Escape should close it
    await page.keyboard.press('Escape');
    await page.waitForTimeout(300);
    const dialogClosedViaEscape = !(await page.locator('text=Choisir une affaire').isVisible());
    console.log(`  Dialog closes via Escape: ${dialogClosedViaEscape}`);

    // ── Step 7: Console errors ───────────────────────────────────────────────
    console.log('\n[Step 7] Console messages:');
    const errors = consoleLogs.filter(m => m.type === 'error' || m.type === 'pageerror');
    const warnings = consoleLogs.filter(m => m.type === 'warning');
    if (errors.length === 0) {
        console.log('  No console errors.');
    } else {
        errors.forEach(e => console.log(`  ERROR: ${e.text}`));
    }
    if (warnings.length === 0) {
        console.log('  No console warnings.');
    } else {
        warnings.forEach(w => console.log(`  WARN:  ${w.text}`));
    }

    await browser.close();
    console.log('\nDone. Screenshots saved to:', screenshotsDir);
})();
