<script>
  import { onMount } from 'svelte';
  import { Chart, registerables } from 'chart.js';
  Chart.register(...registerables);

  let period = 'month'; // 'day' | 'week' | 'month' | 'year'
  let entries = [];
  let goldEntries = [];
  let chartInstance = null;
  let canvasEl;
  let loading = true;
  let timeStr = '';
  let activeDataset = 'both'; // 'money' | 'gold' | 'both'

  function updateCambodiaTime() {
    const now = new Date();
    timeStr = now.toLocaleTimeString('en-US', {
      timeZone: 'Asia/Phnom_Penh',
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  }

  async function loadData() {
    try {
      const [rEntries, rGold] = await Promise.all([
        fetch('/api/entries').then(r => r.json()).catch(() => []),
        fetch('/api/gold').then(r => r.json()).catch(() => [])
      ]);
      entries = Array.isArray(rEntries) ? rEntries : [];
      goldEntries = Array.isArray(rGold) ? rGold : [];
    } catch (e) {
      entries = [];
      goldEntries = [];
    } finally {
      loading = false;
      renderChart();
    }
  }

  // Reactive Stats
  $: totalMoney = entries.reduce((s, e) => s + Number(e.amount || 0), 0);
  $: totalGoldGrams = goldEntries.reduce((s, e) => s + Number(e.grams || 0), 0);
  $: totalGoldMoney = goldEntries.reduce((s, e) => s + Number(e.money || 0), 0);
  $: avgMoneyPerEntry = entries.length > 0 ? (totalMoney / entries.length).toFixed(2) : '0.00';
  $: totalAssetsValue = totalMoney + totalGoldMoney;

  function getWeekString(dateStr) {
    const d = new Date(dateStr + 'T00:00:00');
    const start = new Date(d.getFullYear(), 0, 1);
    const days = Math.floor((d - start) / (24 * 60 * 60 * 1000));
    const weekNumber = Math.ceil((days + start.getDay() + 1) / 7);
    return `${d.getFullYear()}-W${String(weekNumber).padStart(2, '0')}`;
  }

  function aggregateData() {
    const moneyGroup = {};
    const goldGroup = {};
    const allKeysSet = new Set();

    entries.forEach(e => {
      let key = e.date;
      if (period === 'week') key = getWeekString(e.date);
      else if (period === 'month') key = e.date.slice(0, 7);
      else if (period === 'year') key = e.date.slice(0, 4);
      allKeysSet.add(key);
      moneyGroup[key] = (moneyGroup[key] || 0) + Number(e.amount || 0);
    });

    goldEntries.forEach(g => {
      let key = g.date;
      if (period === 'week') key = getWeekString(g.date);
      else if (period === 'month') key = g.date.slice(0, 7);
      else if (period === 'year') key = g.date.slice(0, 4);
      allKeysSet.add(key);
      goldGroup[key] = (goldGroup[key] || 0) + Number(g.grams || 0);
    });

    // Sort labels chronologically
    const sortedKeys = Array.from(allKeysSet).sort();
    const displayKeys = period === 'day' ? sortedKeys.slice(-14) : sortedKeys.slice(-12);

    const moneyData = displayKeys.map(k => moneyGroup[k] || 0);
    const goldData = displayKeys.map(k => goldGroup[k] || 0);

    return { labels: displayKeys, moneyData, goldData };
  }

  function renderChart() {
    if (!canvasEl) return;
    const { labels, moneyData, goldData } = aggregateData();

    if (chartInstance) {
      chartInstance.destroy();
    }

    const ctx = canvasEl.getContext('2d');
    const gradientMoney = ctx.createLinearGradient(0, 0, 0, 300);
    gradientMoney.addColorStop(0, 'rgba(56, 189, 248, 0.85)');
    gradientMoney.addColorStop(1, 'rgba(14, 165, 233, 0.15)');

    const gradientGold = ctx.createLinearGradient(0, 0, 0, 300);
    gradientGold.addColorStop(0, 'rgba(251, 191, 36, 0.85)');
    gradientGold.addColorStop(1, 'rgba(217, 119, 6, 0.15)');

    const datasets = [];

    if (activeDataset === 'money' || activeDataset === 'both') {
      datasets.push({
        label: 'Money Saved ($ USD)',
        data: moneyData,
        backgroundColor: gradientMoney,
        borderColor: '#38bdf8',
        borderWidth: 2,
        borderRadius: 8,
        yAxisID: 'y'
      });
    }

    if (activeDataset === 'gold' || activeDataset === 'both') {
      datasets.push({
        label: 'Gold Saved (Grams)',
        data: goldData,
        backgroundColor: gradientGold,
        borderColor: '#fbbf24',
        borderWidth: 2,
        borderRadius: 8,
        yAxisID: 'y1'
      });
    }

    chartInstance = new Chart(canvasEl, {
      type: 'bar',
      data: {
        labels: labels.length > 0 ? labels : ['No Data'],
        datasets: datasets.length > 0 ? datasets : [{ label: 'No Data', data: [0], backgroundColor: 'rgba(255,255,255,0.05)' }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          mode: 'index',
          intersect: false
        },
        plugins: {
          legend: {
            display: true,
            position: 'top',
            labels: {
              color: 'rgba(255, 255, 255, 0.75)',
              font: { family: 'Plus Jakarta Sans', size: 12, weight: '600' },
              usePointStyle: true,
              pointStyle: 'circle'
            }
          },
          tooltip: {
            backgroundColor: 'rgba(10, 20, 42, 0.95)',
            borderColor: 'rgba(255, 255, 255, 0.2)',
            borderWidth: 1,
            titleColor: '#fff',
            titleFont: { size: 13, weight: '700' },
            bodyFont: { size: 12 },
            padding: 12,
            boxPadding: 6,
            usePointStyle: true,
            callbacks: {
              label: (context) => {
                if (context.dataset.yAxisID === 'y') {
                  return ` Money: $${Number(context.raw).toLocaleString('en-US', { minimumFractionDigits: 2 })}`;
                } else {
                  return ` Gold: ${Number(context.raw).toFixed(3)}g (${(context.raw / 3.75).toFixed(2)} Chi)`;
                }
              }
            }
          }
        },
        scales: {
          x: {
            grid: { display: false },
            ticks: {
              color: 'rgba(255, 255, 255, 0.55)',
              font: { family: 'Plus Jakarta Sans', size: 11, weight: '600' }
            }
          },
          y: {
            type: 'linear',
            display: activeDataset !== 'gold',
            position: 'left',
            beginAtZero: true,
            grid: { color: 'rgba(255, 255, 255, 0.06)' },
            ticks: {
              color: '#38bdf8',
              font: { family: 'Plus Jakarta Sans', size: 11 },
              callback: (v) => `$${v}`
            }
          },
          y1: {
            type: 'linear',
            display: activeDataset !== 'money',
            position: 'right',
            beginAtZero: true,
            grid: { display: false },
            ticks: {
              color: '#fbbf24',
              font: { family: 'Plus Jakarta Sans', size: 11 },
              callback: (v) => `${v}g`
            }
          }
        }
      }
    });
  }

  function setPeriod(p) {
    period = p;
    renderChart();
  }

  function setDataset(d) {
    activeDataset = d;
    renderChart();
  }

  onMount(() => {
    updateCambodiaTime();
    const timer = setInterval(updateCambodiaTime, 1000);
    loadData();
    return () => clearInterval(timer);
  });
</script>

<!-- ══ AMBIENT BACKGROUND SCENE (Voyavera Style) ══ -->
<div class="ambient-scene">
  <div class="ambient-orb ambient-orb-1"></div>
  <div class="ambient-orb ambient-orb-2"></div>
  <div class="ambient-orb ambient-orb-gold"></div>
</div>

<main class="app-container">
  <!-- ══ TOP NAVIGATION (Voyavera Pill Header) ══ -->
  <header class="top-nav glass-panel">
    <div class="brand">
      <div class="brand-logo analytics-glow">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
        </svg>
      </div>
      <span class="brand-name">KHIM X KAIYA</span>
    </div>

    <nav class="nav-pills">
      <a href="index.html" class="nav-pill">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
          <circle cx="12" cy="12" r="10"/><path d="M16 8h-6a2 2 0 1 0 0 4h4a2 2 0 1 1 0 4H8"/><path d="M12 18V6"/>
        </svg>
        <span>Savings</span>
      </a>
      <a href="gold.html" class="nav-pill">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
          <path d="M6 3h12l4 6-10 12L2 9l4-6z"/><path d="M11 3 8 9l4 12 4-12-3-6"/>
        </svg>
        <span>Gold Vault</span>
      </a>
      <a href="analysis.html" class="nav-pill active analytics-pill">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
          <path d="M18 20V10M12 20V4M6 20v-6"/>
        </svg>
        <span>Analytics</span>
      </a>
    </nav>

    <div class="live-clock-badge">
      <span class="live-dot purple"></span>
      <div class="clock-col">
        <span class="clock-digits">{timeStr || '00:00:00'}</span>
        <span class="clock-label">🇰🇭 Live Analytics</span>
      </div>
    </div>
  </header>

  <!-- ══ STATS SUMMARY DECK ══ -->
  <section class="stats-deck">
    <!-- Total Combined Assets -->
    <div class="stat-card glass-panel highlight-blue">
      <div class="sc-top">
        <span class="sc-label">TOTAL CAPITAL VALUE</span>
        <span class="sc-badge">USD + Gold</span>
      </div>
      <div class="sc-amount">${totalAssetsValue.toLocaleString('en-US', { minimumFractionDigits: 2 })}</div>
      <div class="sc-sub">Combined liquid savings & gold holdings</div>
    </div>

    <!-- Total Money Saved -->
    <div class="stat-card glass-panel">
      <div class="sc-top">
        <span class="sc-label">LIQUID SAVINGS</span>
        <span class="sc-icon">💵</span>
      </div>
      <div class="sc-amount">${totalMoney.toLocaleString('en-US', { minimumFractionDigits: 2 })}</div>
      <div class="sc-sub">{entries.length} Deposits • Avg: ${avgMoneyPerEntry}</div>
    </div>

    <!-- Total Gold Saved -->
    <div class="stat-card glass-panel highlight-gold">
      <div class="sc-top">
        <span class="sc-label">TOTAL GOLD RESERVE</span>
        <span class="sc-icon">🥇</span>
      </div>
      <div class="sc-amount">{(totalGoldGrams / 37.5).toFixed(2)} Dom</div>
      <div class="sc-sub">{totalGoldGrams.toFixed(3)}g = {(totalGoldGrams / 3.75).toFixed(2)} Chi</div>
    </div>
  </section>

  <!-- ══ PERIOD & DATASET SELECTION BAR ══ -->
  <section class="controls-bar glass-panel">
    <!-- Period Tabs (Day / Week / Month / Year) -->
    <div class="filter-group">
      <span class="fg-label">TIMEFRAME</span>
      <div class="pill-tabs">
        <button class="pt-btn" class:active={period === 'day'} on:click={() => setPeriod('day')}>
          Day
        </button>
        <button class="pt-btn" class:active={period === 'week'} on:click={() => setPeriod('week')}>
          Week
        </button>
        <button class="pt-btn" class:active={period === 'month'} on:click={() => setPeriod('month')}>
          Month
        </button>
        <button class="pt-btn" class:active={period === 'year'} on:click={() => setPeriod('year')}>
          Year
        </button>
      </div>
    </div>

    <!-- Dataset Filter (Both / Money / Gold) -->
    <div class="filter-group">
      <span class="fg-label">DATASET VIEW</span>
      <div class="pill-tabs">
        <button class="pt-btn" class:active={activeDataset === 'both'} on:click={() => setDataset('both')}>
          All Assets
        </button>
        <button class="pt-btn" class:active={activeDataset === 'money'} on:click={() => setDataset('money')}>
          💵 Money
        </button>
        <button class="pt-btn" class:active={activeDataset === 'gold'} on:click={() => setDataset('gold')}>
          🥇 Gold
        </button>
      </div>
    </div>
  </section>

  <!-- ══ CHART CANVAS CARD ══ -->
  <section class="chart-section glass-panel">
    <div class="cs-header">
      <div>
        <h3 class="cs-title">Accumulation Trajectory</h3>
        <p class="cs-sub">
          Aggregated by {period === 'day' ? 'daily deposits' : period === 'week' ? 'weekly intervals' : period === 'month' ? 'monthly totals' : 'annual totals'}
        </p>
      </div>
      <div class="cs-legend-pill">
        <span class="legend-dot blue"></span> <span>Money (USD)</span>
        <span class="legend-dot gold"></span> <span>Gold (Grams)</span>
      </div>
    </div>

    <div class="canvas-container">
      <canvas bind:this={canvasEl}></canvas>
    </div>
  </section>

  <!-- ══ BREAKDOWN LEDGERS (Side-by-side Grid) ══ -->
  <section class="ledgers-grid">
    <!-- Money Ledger Preview -->
    <div class="ledger-column glass-panel">
      <div class="lc-header">
        <div class="lc-title-wrap">
          <span class="lc-icon">💵</span>
          <h4 class="lc-title">Recent Money Deposits</h4>
        </div>
        <a href="index.html" class="lc-link">View All →</a>
      </div>

      <div class="lc-list">
        {#if entries.length === 0}
          <div class="lc-empty">No money transactions yet</div>
        {:else}
          {#each entries.slice(0, 6) as item}
            <div class="lc-item">
              <div class="lci-left">
                <span class="lci-date">{item.date}</span>
                <span class="lci-note">{item.note || 'General Savings'}</span>
              </div>
              <span class="lci-amount money">+${Number(item.amount).toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- Gold Ledger Preview -->
    <div class="ledger-column glass-panel">
      <div class="lc-header">
        <div class="lc-title-wrap">
          <span class="lc-icon">🥇</span>
          <h4 class="lc-title">Recent Gold Holdings</h4>
        </div>
        <a href="gold.html" class="lc-link">View All →</a>
      </div>

      <div class="lc-list">
        {#if goldEntries.length === 0}
          <div class="lc-empty">No gold entries yet</div>
        {:else}
          {#each goldEntries.slice(0, 6) as item}
            <div class="lc-item">
              <div class="lci-left">
                <span class="lci-date">{item.date}</span>
                <span class="lci-note">{item.qty} {item.unit} ({Number(item.grams).toFixed(3)}g)</span>
              </div>
              <span class="lci-amount gold">
                {#if item.money > 0}
                  ${Number(item.money).toLocaleString('en-US', { minimumFractionDigits: 2 })}
                {:else}
                  {item.qty} {item.unit}
                {/if}
              </span>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </section>
</main>

<style>
  /* ══ NAVIGATION ══ */
  .top-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 24px;
    border-radius: var(--radius-pill);
    margin-bottom: 28px;
    background: rgba(10, 20, 42, 0.65);
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .brand-logo {
    width: 38px;
    height: 38px;
    border-radius: 12px;
    background: linear-gradient(135deg, #f43f5e, #c084fc);
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    box-shadow: 0 4px 16px rgba(244, 63, 94, 0.45);
  }
  .brand-name {
    font-size: 19px;
    font-weight: 800;
    letter-spacing: -0.2px;
    background: linear-gradient(135deg, #ffffff 40%, #e9d5ff 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .nav-pills {
    display: flex;
    gap: 6px;
    background: rgba(255, 255, 255, 0.06);
    padding: 4px;
    border-radius: var(--radius-pill);
    border: 1px solid rgba(255, 255, 255, 0.08);
  }
  .nav-pill {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 18px;
    border-radius: var(--radius-pill);
    font-size: 13.5px;
    font-weight: 600;
    color: var(--text-secondary);
    text-decoration: none;
    transition: all 0.25s ease;
  }
  .nav-pill:hover {
    color: #fff;
    background: rgba(255, 255, 255, 0.08);
  }
  .nav-pill.active.analytics-pill {
    background: rgba(139, 92, 246, 0.2);
    color: #e9d5ff;
    border: 1px solid rgba(139, 92, 246, 0.45);
    box-shadow: 0 4px 18px rgba(139, 92, 246, 0.35);
  }

  .live-clock-badge {
    display: flex;
    align-items: center;
    gap: 10px;
    background: rgba(255, 255, 255, 0.06);
    padding: 7px 14px;
    border-radius: var(--radius-pill);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  .live-dot.purple {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #a855f7;
    box-shadow: 0 0 10px #a855f7;
    animation: livePulse 2s infinite;
  }
  @keyframes livePulse {
    0%, 100% { opacity: 1; transform: scale(1); }
    50% { opacity: 0.4; transform: scale(0.85); }
  }
  .clock-col {
    display: flex;
    flex-direction: column;
    line-height: 1.1;
  }
  .clock-digits {
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.5px;
  }
  .clock-label {
    font-size: 9.5px;
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  /* ══ STATS DECK ══ */
  .stats-deck {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 18px;
    margin-bottom: 26px;
  }
  .stat-card {
    padding: 24px;
    border-radius: 28px;
  }
  .stat-card.highlight-blue {
    background: linear-gradient(135deg, rgba(37, 99, 235, 0.35) 0%, rgba(14, 165, 233, 0.2) 100%);
    border-color: rgba(56, 189, 248, 0.35);
  }
  .stat-card.highlight-gold {
    background: linear-gradient(135deg, rgba(60, 35, 15, 0.45) 0%, rgba(30, 20, 10, 0.35) 100%);
    border-color: rgba(251, 191, 36, 0.35);
  }
  .sc-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }
  .sc-label {
    font-size: 11px;
    font-weight: 800;
    letter-spacing: 1.2px;
    color: var(--text-tertiary);
  }
  .sc-badge {
    background: rgba(255, 255, 255, 0.12);
    padding: 2px 8px;
    border-radius: var(--radius-pill);
    font-size: 10.5px;
    font-weight: 700;
    color: #38bdf8;
  }
  .sc-icon {
    font-size: 18px;
  }
  .sc-amount {
    font-size: 32px;
    font-weight: 800;
    letter-spacing: -1px;
    color: #fff;
    margin-bottom: 6px;
  }
  .sc-sub {
    font-size: 12.5px;
    color: var(--text-secondary);
  }

  /* ══ CONTROLS BAR ══ */
  .controls-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    border-radius: 24px;
    margin-bottom: 24px;
  }
  .filter-group {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .fg-label {
    font-size: 11px;
    font-weight: 800;
    letter-spacing: 1px;
    color: var(--text-tertiary);
  }
  .pill-tabs {
    display: flex;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: var(--radius-pill);
    padding: 3px;
    gap: 3px;
  }
  .pt-btn {
    padding: 7px 16px;
    border-radius: var(--radius-pill);
    background: none;
    border: none;
    font-size: 13px;
    font-weight: 600;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s;
  }
  .pt-btn:hover {
    color: #fff;
  }
  .pt-btn.active {
    background: rgba(255, 255, 255, 0.18);
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.25);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  /* ══ CHART SECTION ══ */
  .chart-section {
    padding: 28px;
    border-radius: 32px;
    margin-bottom: 28px;
  }
  .cs-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
  }
  .cs-title {
    font-size: 20px;
    font-weight: 800;
    color: #fff;
    margin-bottom: 4px;
  }
  .cs-sub {
    font-size: 13px;
    color: var(--text-tertiary);
  }
  .cs-legend-pill {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 12px;
    color: var(--text-secondary);
    background: rgba(255, 255, 255, 0.06);
    padding: 6px 14px;
    border-radius: var(--radius-pill);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  .legend-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }
  .legend-dot.blue { background: #38bdf8; }
  .legend-dot.gold { background: #fbbf24; }

  .canvas-container {
    height: 340px;
    position: relative;
  }

  /* ══ LEDGERS GRID ══ */
  .ledgers-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }
  .ledger-column {
    padding: 24px;
    border-radius: 28px;
  }
  .lc-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  }
  .lc-title-wrap {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .lc-icon {
    font-size: 18px;
  }
  .lc-title {
    font-size: 15px;
    font-weight: 700;
    color: #fff;
  }
  .lc-link {
    font-size: 12.5px;
    font-weight: 600;
    color: #38bdf8;
    text-decoration: none;
  }

  .lc-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .lc-empty {
    text-align: center;
    padding: 24px;
    font-size: 13px;
    color: var(--text-tertiary);
  }
  .lc-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 14px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 14px;
  }
  .lci-left {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .lci-date {
    font-size: 13px;
    font-weight: 700;
    color: #fff;
  }
  .lci-note {
    font-size: 11.5px;
    color: var(--text-tertiary);
  }
  .lci-amount {
    font-size: 14.5px;
    font-weight: 800;
  }
  .lci-amount.money { color: #34d399; }
  .lci-amount.gold { color: #fbbf24; }

  /* ══ RESPONSIVE ADAPTATIONS (Desktop, Tablet, Mobile) ══ */
  @media (max-width: 860px) {
    .stats-deck {
      grid-template-columns: 1fr;
      gap: 14px;
    }
    .ledgers-grid {
      grid-template-columns: 1fr;
      gap: 16px;
    }
    .controls-bar {
      flex-direction: column;
      gap: 14px;
      align-items: stretch;
      padding: 16px;
    }
    .filter-group {
      justify-content: space-between;
    }
    .top-nav {
      flex-wrap: wrap;
      gap: 12px;
      border-radius: 24px;
      padding: 14px 18px;
    }
    .chart-section {
      padding: 22px 18px;
    }
    .canvas-container {
      height: 300px;
    }
  }

  @media (max-width: 640px) {
    .top-nav {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;
      padding: 14px;
      border-radius: 20px;
      margin-bottom: 20px;
    }
    .brand {
      justify-content: center;
    }
    .nav-pills {
      width: 100%;
      justify-content: center;
      gap: 4px;
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
      padding: 4px;
    }
    .nav-pill {
      padding: 8px 12px;
      font-size: 13px;
      flex: 1;
      justify-content: center;
    }
    .live-clock-badge {
      align-self: center;
      padding: 5px 12px;
    }
    .stats-deck {
      margin-bottom: 20px;
    }
    .stat-card {
      padding: 18px 16px;
      border-radius: 20px;
    }
    .sc-amount {
      font-size: clamp(26px, 7vw, 34px);
    }
    .controls-bar {
      padding: 14px 12px;
      border-radius: 20px;
      margin-bottom: 20px;
      gap: 12px;
    }
    .filter-group {
      flex-direction: column;
      align-items: stretch;
      gap: 8px;
    }
    .fg-label {
      font-size: 11px;
    }
    .pill-tabs {
      width: 100%;
      display: flex;
    }
    .pt-btn {
      flex: 1;
      text-align: center;
      justify-content: center;
      padding: 8px 6px;
      font-size: 12px;
    }
    .chart-section {
      padding: 16px 14px;
      border-radius: 22px;
      margin-bottom: 20px;
    }
    .cs-header {
      flex-direction: column;
      gap: 10px;
      margin-bottom: 16px;
    }
    .cs-title {
      font-size: 18px;
    }
    .canvas-container {
      height: 240px;
    }
    .ledger-column {
      padding: 18px 14px;
      border-radius: 20px;
    }
    .lc-item {
      padding: 10px 0;
    }
  }

  @media (max-width: 420px) {
    .nav-pill span {
      font-size: 12px;
    }
    .sc-amount {
      font-size: 24px;
    }
    .pt-btn {
      font-size: 11px;
      padding: 7px 4px;
    }
  }
</style>
