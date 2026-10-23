<script>
  import { onMount } from 'svelte';

  let entries = [];
  let datePicker = new Date().toISOString().slice(0, 10);
  let amountInput = '';
  let noteInput = '';
  let total = 0;
  let weather = null;
  let timeStr = '';
  let dateStr = '';
  let editingDate = null;
  let showModal = false;
  let loading = true;
  let searchQuery = '';
  let alert = { show: false, msg: '', type: 'success', icon: '✨' };
  let alertTimer = null;

  // Preset quick amounts
  const quickAmounts = [5, 10, 20, 50, 100, 500];

  function triggerAlert(msg, type = 'success', icon = '✨') {
    if (alertTimer) clearTimeout(alertTimer);
    alert = { show: true, msg, type, icon };
    alertTimer = setTimeout(() => {
      alert = { ...alert, show: false };
    }, 2800);
  }

  function updateCambodiaTime() {
    const now = new Date();
    timeStr = now.toLocaleTimeString('en-US', {
      timeZone: 'Asia/Phnom_Penh',
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
    dateStr = now.toLocaleDateString('en-US', {
      timeZone: 'Asia/Phnom_Penh',
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }

  async function fetchWeather() {
    try {
      const res = await fetch('/api/weather');
      const data = await res.json();
      if (data && data.temp !== undefined) {
        weather = data;
      }
    } catch (e) {
      weather = { temp: 28, description: 'Partly cloudy', humidity: 75, wind_speed: 10, icon: '⛅' };
    }
  }

  async function fetchEntries() {
    try {
      const res = await fetch('/api/entries');
      const data = await res.json();
      if (Array.isArray(data)) {
        entries = data;
        total = entries.reduce((sum, item) => sum + Number(item.amount || 0), 0);
      } else {
        entries = [];
        total = 0;
      }
    } catch (e) {
      entries = [];
      total = 0;
    } finally {
      loading = false;
    }
  }

  async function handleSave() {
    const num = parseFloat(amountInput);
    if (!datePicker || isNaN(num) || num <= 0) {
      triggerAlert('Please enter a valid amount', 'error', '⚠️');
      return;
    }

    try {
      const res = await fetch('/api/entry', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          date: datePicker,
          amount: num,
          note: noteInput.trim()
        })
      });

      if (res.ok) {
        triggerAlert(
          editingDate ? 'Entry updated successfully!' : `Saved $${num.toFixed(2)} to Vault!`,
          'success',
          '✓'
        );
        closeModal();
        await fetchEntries();
      } else {
        triggerAlert('Failed to save entry', 'error', '✕');
      }
    } catch (e) {
      triggerAlert('Network error saving entry', 'error', '✕');
    }
  }

  function openAddModal() {
    editingDate = null;
    datePicker = new Date().toISOString().slice(0, 10);
    amountInput = '';
    noteInput = '';
    showModal = true;
  }

  function openEditModal(entry) {
    editingDate = entry.date;
    datePicker = entry.date;
    amountInput = entry.amount;
    noteInput = entry.note || '';
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    editingDate = null;
  }

  async function deleteEntry(date) {
    if (!confirm(`Delete entry for ${date}?`)) return;
    try {
      const res = await fetch(`/api/entry/${date}`, { method: 'DELETE' });
      if (res.ok) {
        triggerAlert(`Deleted entry for ${date}`, 'error', '🗑️');
        await fetchEntries();
      }
    } catch (e) {
      triggerAlert('Could not delete entry', 'error', '✕');
    }
  }

  function setQuickAmount(val) {
    const current = parseFloat(amountInput) || 0;
    amountInput = (current + val).toString();
  }

  $: filteredEntries = entries.filter(e => {
    if (!searchQuery.trim()) return true;
    const q = searchQuery.toLowerCase();
    return (e.date && e.date.toLowerCase().includes(q)) ||
           (e.note && e.note.toLowerCase().includes(q)) ||
           (e.amount && e.amount.toString().includes(q));
  });

  onMount(() => {
    updateCambodiaTime();
    const clockTimer = setInterval(updateCambodiaTime, 1000);
    fetchWeather();
    fetchEntries();
    const weatherTimer = setInterval(fetchWeather, 5 * 60 * 1000);

    return () => {
      clearInterval(clockTimer);
      clearInterval(weatherTimer);
    };
  });
</script>

<!-- ══ AMBIENT BACKGROUND SCENE (Voyavera Vibe) ══ -->
<div class="ambient-scene">
  <div class="ambient-orb ambient-orb-1"></div>
  <div class="ambient-orb ambient-orb-2"></div>
  <div class="ambient-orb ambient-orb-3"></div>
</div>

<!-- ══ DYNAMIC ISLAND NOTIFICATION ALERT ══ -->
{#if alert.show}
  <div class="dynamic-island-alert {alert.type}">
    <span class="dia-icon">{alert.icon}</span>
    <span class="dia-text">{alert.msg}</span>
    <div class="dia-glow-bar"></div>
  </div>
{/if}

<main class="app-container">
  <!-- ══ TOP NAVIGATION (Voyavera Style Pill Header) ══ -->
  <header class="top-nav glass-panel">
    <div class="brand">
      <div class="brand-logo">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
        </svg>
      </div>
      <span class="brand-name">KHIM X KAIYA</span>
    </div>

    <!-- Center Nav Pills -->
    <nav class="nav-pills">
      <a href="index.html" class="nav-pill active">
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
      <a href="analysis.html" class="nav-pill">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
          <path d="M18 20V10M12 20V4M6 20v-6"/>
        </svg>
        <span>Analytics</span>
      </a>
    </nav>

    <!-- Right Live Cambodia Badge -->
    <div class="live-clock-badge">
      <span class="live-dot"></span>
      <div class="clock-col">
        <span class="clock-digits">{timeStr || '00:00:00'}</span>
        <span class="clock-label">🇰🇭 Cambodia</span>
      </div>
    </div>
  </header>

  <!-- ══ HERO BANNER (Atmospheric Mountain & Sky Glass Card) ══ -->
  <section class="hero-card glass-panel">
    <div class="hero-grid">
      <!-- Left Content -->
      <div class="hero-left">
        <div class="eyebrow-tag">
          <span class="tag-sparkle">✦</span>
          <span>DISCOVER • BUILD • ACCUMULATE</span>
        </div>
        <h1 class="hero-headline">
          The Wealth <br />
          <span class="gradient-text">Awaits You</span>
        </h1>
        <p class="hero-subtitle">
          Track daily savings in USD and Cambodian gold units with effortless real-time analytics.
        </p>

        <!-- CTAs -->
        <div class="hero-actions">
          <button class="btn-voyavera-primary" on:click={openAddModal}>
            <span>Add Savings</span>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </button>
          <a href="gold.html" class="btn-voyavera-glass">
            <span>Explore Gold</span>
            <span class="gold-sparkle">🥇</span>
          </a>
        </div>
      </div>

      <!-- Right Glass Card Widget: Live Total & Cambodia Weather -->
      <div class="hero-right">
        <!-- Floating Weather Glass Capsule (Voyavera Switzerland style) -->
        {#if weather}
          <div class="weather-floating-card">
            <div class="wfc-header">
              <div class="wfc-icon-wrap">{weather.icon}</div>
              <div>
                <div class="wfc-city">Phnom Penh, KH</div>
                <div class="wfc-condition">{weather.description}</div>
              </div>
            </div>
            <div class="wfc-body">
              <span class="wfc-temp">{weather.temp}°C</span>
              <div class="wfc-meta">
                <span>💧 {weather.humidity}% hum</span>
                <span>💨 {weather.wind_speed} km/h</span>
              </div>
            </div>
            <div class="wfc-date">{dateStr}</div>
          </div>
        {/if}

        <!-- Total Balance Floating Card -->
        <div class="balance-card">
          <div class="bc-top">
            <span class="bc-label">TOTAL CAPITAL SAVED</span>
            <span class="bc-badge">{entries.length} Deposits</span>
          </div>
          <div class="bc-amount">
            <span class="bc-currency">$</span>
            <span class="bc-number">{total.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
          </div>
          <div class="bc-footer">
            <span class="trend-indicator">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/>
              </svg>
              Secured in SQLite DB
            </span>
            <span class="bc-action-link" on:click={openAddModal}>+ Deposit Now</span>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- ══ ENTRIES & ACTIVITY SECTION ══ -->
  <section class="activity-section">
    <div class="activity-header">
      <div>
        <h2 class="section-heading">Savings Journal</h2>
        <p class="section-sub">Chronological ledger of all daily deposits</p>
      </div>

      <!-- Search & Add Actions -->
      <div class="activity-tools">
        <div class="search-input-wrap">
          <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
            <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            type="text"
            placeholder="Search dates, notes..."
            bind:value={searchQuery}
            class="search-input"
          />
          {#if searchQuery}
            <button class="clear-search" on:click={() => searchQuery = ''}>✕</button>
          {/if}
        </div>
        <button class="btn-add-circle" on:click={openAddModal} title="New Deposit">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Entries Table / Grid -->
    {#if loading}
      <div class="skeleton-wrap">
        {#each [1, 2, 3, 4] as _}
          <div class="skeleton-card"></div>
        {/each}
      </div>
    {:else if filteredEntries.length === 0}
      <div class="empty-glass-card glass-panel">
        <div class="empty-icon-halo">
          <svg width="44" height="44" viewBox="0 0 24 24" fill="none" stroke="rgba(56, 189, 248, 0.7)" stroke-width="1.8">
            <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
          </svg>
        </div>
        <h3 class="empty-title">{searchQuery ? 'No matching deposits' : 'No savings deposited yet'}</h3>
        <p class="empty-desc">
          {searchQuery ? 'Try adjusting your search criteria' : 'Click the button below to record your first saving entry'}
        </p>
        <button class="btn-voyavera-primary sm" on:click={openAddModal}>
          <span>+ Deposit Money</span>
        </button>
      </div>
    {:else}
      <div class="entries-grid">
        {#each filteredEntries as entry, index}
          <div class="entry-row glass-panel" style="animation-delay: {index * 35}ms">
            <!-- Left Info -->
            <div class="er-left">
              <div class="er-date-box">
                <span class="er-day">{entry.date.split('-')[2] || '01'}</span>
                <span class="er-month">{new Date(entry.date).toLocaleDateString('en-US', { month: 'short' })}</span>
              </div>
              <div class="er-details">
                <div class="er-full-date">{entry.date}</div>
                <div class="er-note">{entry.note || 'General Savings'}</div>
              </div>
            </div>

            <!-- Right Amount & Controls -->
            <div class="er-right">
              <div class="er-amount-wrap">
                <span class="er-plus">+</span>
                <span class="er-amount">${Number(entry.amount).toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
              </div>
              <div class="er-actions">
                <button class="action-btn edit" on:click={() => openEditModal(entry)} title="Edit">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                </button>
                <button class="action-btn delete" on:click={() => deleteEntry(entry.date)} title="Delete">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
                    <polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </section>
</main>

<!-- ══ ADD / EDIT MODAL (Voyavera Frosted Glass Sheet) ══ -->
{#if showModal}
  <div class="modal-backdrop" on:click={closeModal}>
    <div class="modal-dialog glass-panel" on:click|stopPropagation>
      <div class="modal-top">
        <div class="modal-title-group">
          <div class="modal-tag">SAVINGS TRANSACTION</div>
          <h3 class="modal-heading">{editingDate ? 'Update Savings' : 'Deposit Savings'}</h3>
        </div>
        <button class="modal-close" on:click={closeModal}>✕</button>
      </div>

      <div class="modal-body">
        <!-- Date input -->
        <div class="form-field">
          <label class="ff-label">Transaction Date</label>
          <div class="ff-input-wrap">
            <input type="date" class="ff-input" bind:value={datePicker} />
          </div>
        </div>

        <!-- Quick Amount Pills -->
        <div class="form-field">
          <label class="ff-label">Quick Add (USD)</label>
          <div class="quick-pill-grid">
            {#each quickAmounts as qa}
              <button
                type="button"
                class="quick-pill"
                on:click={() => setQuickAmount(qa)}
              >
                +${qa}
              </button>
            {/each}
          </div>
        </div>

        <!-- Amount Input -->
        <div class="form-field">
          <label class="ff-label">Amount (USD)</label>
          <div class="ff-input-wrap prefix-wrap">
            <span class="ff-prefix">$</span>
            <input
              type="number"
              class="ff-input prefix"
              placeholder="0.00"
              step="0.01"
              min="0.01"
              bind:value={amountInput}
              autofocus
            />
          </div>
        </div>

        <!-- Note Input -->
        <div class="form-field">
          <label class="ff-label">Note / Purpose (Optional)</label>
          <input
            type="text"
            class="ff-input"
            placeholder="e.g. Salary, Dividend, Coffee savings..."
            bind:value={noteInput}
          />
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-cancel" on:click={closeModal}>Cancel</button>
        <button class="btn-voyavera-primary" on:click={handleSave}>
          <span>{editingDate ? 'Save Changes' : 'Confirm Deposit'}</span>
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
{/if}

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
    background: linear-gradient(135deg, #f43f5e, #fb7185);
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
    background: linear-gradient(135deg, #ffffff 40%, #fecdd3 100%);
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
  .nav-pill.active {
    background: rgba(255, 255, 255, 0.16);
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.25);
    box-shadow: 0 4px 15px rgba(14, 165, 233, 0.25);
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
  .live-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #10b981;
    box-shadow: 0 0 8px #10b981;
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

  /* ══ HERO BANNER ══ */
  .hero-card {
    padding: 38px 42px;
    margin-bottom: 34px;
    border-radius: 36px;
    background: linear-gradient(135deg, rgba(14, 28, 58, 0.82) 0%, rgba(9, 18, 38, 0.72) 100%);
  }
  .hero-grid {
    display: grid;
    grid-template-columns: 1.15fr 0.85fr;
    gap: 36px;
    align-items: center;
  }

  .eyebrow-tag {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 1.5px;
    color: #38bdf8;
    margin-bottom: 14px;
  }
  .tag-sparkle {
    font-size: 12px;
  }

  .hero-headline {
    font-family: var(--font-display);
    font-size: 48px;
    font-weight: 800;
    line-height: 1.08;
    letter-spacing: -1.8px;
    margin-bottom: 16px;
    color: #fff;
  }
  .gradient-text {
    background: linear-gradient(135deg, #38bdf8 0%, #818cf8 50%, #c084fc 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .hero-subtitle {
    font-size: 15px;
    color: var(--text-secondary);
    line-height: 1.5;
    max-width: 440px;
    margin-bottom: 28px;
  }

  .hero-actions {
    display: flex;
    gap: 14px;
    align-items: center;
  }

  .btn-voyavera-primary {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    background: #ffffff;
    color: #070e1e;
    padding: 13px 26px;
    border-radius: var(--radius-pill);
    font-size: 14.5px;
    font-weight: 700;
    border: none;
    cursor: pointer;
    text-decoration: none;
    box-shadow: 0 10px 25px -5px rgba(255, 255, 255, 0.35);
    transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  }
  .btn-voyavera-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 14px 30px -5px rgba(255, 255, 255, 0.5);
  }
  .btn-voyavera-primary:active {
    transform: translateY(0) scale(0.98);
  }
  .btn-voyavera-primary.sm {
    padding: 10px 20px;
    font-size: 13.5px;
  }

  .btn-voyavera-glass {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: rgba(255, 255, 255, 0.08);
    color: #ffffff;
    padding: 13px 24px;
    border-radius: var(--radius-pill);
    font-size: 14.5px;
    font-weight: 600;
    border: 1px solid rgba(255, 255, 255, 0.2);
    text-decoration: none;
    backdrop-filter: blur(20px);
    transition: all 0.25s ease;
  }
  .btn-voyavera-glass:hover {
    background: rgba(255, 255, 255, 0.16);
    border-color: rgba(255, 255, 255, 0.35);
    transform: translateY(-2px);
  }

  /* Right Widgets */
  .hero-right {
    display: flex;
    flex-direction: column;
    gap: 18px;
  }

  .weather-floating-card {
    background: rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(24px);
    border: 1px solid rgba(255, 255, 255, 0.18);
    border-radius: 24px;
    padding: 18px 22px;
    box-shadow: 0 16px 36px -10px rgba(0, 0, 0, 0.4);
    animation: fadeIn 0.6s ease;
  }
  .wfc-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 10px;
  }
  .wfc-icon-wrap {
    font-size: 28px;
    line-height: 1;
  }
  .wfc-city {
    font-size: 14px;
    font-weight: 700;
    color: #fff;
  }
  .wfc-condition {
    font-size: 12px;
    color: #93c5fd;
    text-transform: capitalize;
  }
  .wfc-body {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 8px;
  }
  .wfc-temp {
    font-size: 26px;
    font-weight: 800;
    color: #fff;
    letter-spacing: -0.5px;
  }
  .wfc-meta {
    font-size: 11.5px;
    color: var(--text-secondary);
    display: flex;
    gap: 10px;
  }
  .wfc-date {
    font-size: 11px;
    color: var(--text-tertiary);
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    padding-top: 6px;
  }

  .balance-card {
    background: linear-gradient(135deg, rgba(37, 99, 235, 0.4) 0%, rgba(14, 165, 233, 0.25) 100%);
    border: 1px solid rgba(56, 189, 248, 0.35);
    border-radius: 24px;
    padding: 22px;
    box-shadow: 0 20px 40px -10px rgba(14, 165, 233, 0.3);
  }
  .bc-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
  }
  .bc-label {
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 1px;
    color: #93c5fd;
  }
  .bc-badge {
    background: rgba(255, 255, 255, 0.15);
    padding: 2px 10px;
    border-radius: var(--radius-pill);
    font-size: 11px;
    font-weight: 600;
  }
  .bc-amount {
    display: flex;
    align-items: baseline;
    gap: 4px;
    margin-bottom: 14px;
  }
  .bc-currency {
    font-size: 24px;
    font-weight: 700;
    color: #38bdf8;
  }
  .bc-number {
    font-size: 38px;
    font-weight: 800;
    letter-spacing: -1.2px;
    color: #ffffff;
  }
  .bc-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 10px;
    border-top: 1px solid rgba(255, 255, 255, 0.12);
  }
  .trend-indicator {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: #6ee7b7;
    font-weight: 600;
  }
  .bc-action-link {
    font-size: 12px;
    font-weight: 700;
    color: #fff;
    cursor: pointer;
    text-decoration: underline;
    transition: opacity 0.2s;
  }
  .bc-action-link:hover {
    opacity: 0.8;
  }

  /* ══ ACTIVITY JOURNAL ══ */
  .activity-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 22px;
  }
  .section-heading {
    font-size: 24px;
    font-weight: 800;
    letter-spacing: -0.6px;
    color: #fff;
    margin-bottom: 4px;
  }
  .section-sub {
    font-size: 13.5px;
    color: var(--text-tertiary);
  }

  .activity-tools {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .search-input-wrap {
    position: relative;
    width: 220px;
  }
  .search-icon {
    position: absolute;
    left: 14px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-tertiary);
  }
  .search-input {
    width: 100%;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: var(--radius-pill);
    padding: 9px 34px;
    font-size: 13px;
    color: #fff;
    outline: none;
    transition: all 0.2s;
  }
  .search-input:focus {
    background: rgba(255, 255, 255, 0.1);
    border-color: #38bdf8;
    width: 250px;
  }
  .clear-search {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    font-size: 12px;
  }

  .btn-add-circle {
    width: 38px;
    height: 38px;
    border-radius: 50%;
    background: linear-gradient(135deg, #0ea5e9, #38bdf8);
    border: none;
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    box-shadow: 0 4px 14px rgba(14, 165, 233, 0.4);
    transition: all 0.2s;
  }
  .btn-add-circle:hover {
    transform: scale(1.08);
  }

  /* Entries Grid */
  .entries-grid {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .entry-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 22px;
    border-radius: 20px;
    animation: slideUp 0.35s ease both;
  }
  @keyframes slideUp {
    from { opacity: 0; transform: translateY(14px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .er-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .er-date-box {
    width: 48px;
    height: 48px;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 14px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    line-height: 1;
  }
  .er-day {
    font-size: 17px;
    font-weight: 800;
    color: #fff;
  }
  .er-month {
    font-size: 10px;
    font-weight: 700;
    color: #38bdf8;
    text-transform: uppercase;
    margin-top: 2px;
  }
  .er-details {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .er-full-date {
    font-size: 14.5px;
    font-weight: 700;
    color: #fff;
  }
  .er-note {
    font-size: 12.5px;
    color: var(--text-secondary);
  }

  .er-right {
    display: flex;
    align-items: center;
    gap: 18px;
  }
  .er-amount-wrap {
    display: flex;
    align-items: baseline;
    gap: 2px;
  }
  .er-plus {
    font-size: 13px;
    font-weight: 700;
    color: #34d399;
  }
  .er-amount {
    font-size: 18px;
    font-weight: 800;
    color: #34d399;
    letter-spacing: -0.3px;
  }
  .er-actions {
    display: flex;
    gap: 6px;
  }
  .action-btn {
    width: 32px;
    height: 32px;
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }
  .action-btn:hover.edit {
    background: rgba(56, 189, 248, 0.2);
    color: #38bdf8;
    border-color: rgba(56, 189, 248, 0.4);
  }
  .action-btn:hover.delete {
    background: rgba(244, 63, 94, 0.2);
    color: #f43f5e;
    border-color: rgba(244, 63, 94, 0.4);
  }

  /* Empty State */
  .empty-glass-card {
    text-align: center;
    padding: 56px 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .empty-icon-halo {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: rgba(56, 189, 248, 0.1);
    border: 1px solid rgba(56, 189, 248, 0.25);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 18px;
    box-shadow: 0 0 30px rgba(56, 189, 248, 0.2);
  }
  .empty-title {
    font-size: 19px;
    font-weight: 700;
    margin-bottom: 6px;
  }
  .empty-desc {
    font-size: 13.5px;
    color: var(--text-tertiary);
    max-width: 340px;
    margin-bottom: 22px;
  }

  /* Skeleton Loading */
  .skeleton-wrap {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .skeleton-card {
    height: 76px;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.05);
    animation: shimmer 1.6s infinite;
  }
  @keyframes shimmer {
    0%, 100% { opacity: 0.4; }
    50% { opacity: 0.8; }
  }

  /* ══ MODAL SHEET ══ */
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(2, 6, 18, 0.72);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    animation: fadeIn 0.25s ease;
  }
  .modal-dialog {
    width: 100%;
    max-width: 480px;
    background: rgba(13, 24, 48, 0.94);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 32px;
    padding: 30px;
    box-shadow: 0 25px 60px -15px rgba(0, 0, 0, 0.8), 0 0 40px rgba(14, 165, 233, 0.2);
    animation: scaleUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  }
  @keyframes scaleUp {
    from { opacity: 0; transform: scale(0.92) translateY(20px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
  }

  .modal-top {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 22px;
  }
  .modal-tag {
    font-size: 10.5px;
    font-weight: 700;
    letter-spacing: 1.2px;
    color: #38bdf8;
    margin-bottom: 4px;
  }
  .modal-heading {
    font-size: 22px;
    font-weight: 800;
    color: #fff;
    letter-spacing: -0.4px;
  }
  .modal-close {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.15);
    color: #fff;
    cursor: pointer;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }
  .modal-close:hover {
    background: rgba(255, 255, 255, 0.18);
  }

  .modal-body {
    display: flex;
    flex-direction: column;
    gap: 18px;
    margin-bottom: 26px;
  }

  .form-field {
    display: flex;
    flex-direction: column;
    gap: 7px;
  }
  .ff-label {
    font-size: 11.5px;
    font-weight: 700;
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.6px;
  }
  .ff-input-wrap {
    position: relative;
  }
  .ff-input {
    width: 100%;
    background: rgba(255, 255, 255, 0.06);
    border: 1.5px solid rgba(255, 255, 255, 0.12);
    border-radius: 14px;
    padding: 13px 16px;
    font-size: 15px;
    color: #fff;
    font-family: inherit;
    outline: none;
    transition: all 0.2s;
  }
  .ff-input:focus {
    border-color: #38bdf8;
    background: rgba(56, 189, 248, 0.08);
    box-shadow: 0 0 0 4px rgba(56, 189, 248, 0.15);
  }
  .prefix-wrap {
    position: relative;
  }
  .ff-prefix {
    position: absolute;
    left: 16px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 18px;
    font-weight: 700;
    color: #38bdf8;
  }
  .ff-input.prefix {
    padding-left: 36px;
    font-size: 18px;
    font-weight: 700;
  }

  .quick-pill-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 6px;
  }
  .quick-pill {
    background: rgba(255, 255, 255, 0.07);
    border: 1px solid rgba(255, 255, 255, 0.14);
    border-radius: 10px;
    padding: 7px 4px;
    font-size: 12px;
    font-weight: 700;
    color: #fff;
    cursor: pointer;
    transition: all 0.15s;
    text-align: center;
  }
  .quick-pill:hover {
    background: rgba(56, 189, 248, 0.2);
    border-color: #38bdf8;
  }
  .quick-pill:active {
    transform: scale(0.94);
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
  .btn-cancel {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.14);
    color: var(--text-secondary);
    padding: 12px 20px;
    border-radius: var(--radius-pill);
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }
  .btn-cancel:hover {
    background: rgba(255, 255, 255, 0.14);
    color: #fff;
  }

  /* ══ RESPONSIVE ADAPTATIONS (Desktop, Tablet, Mobile) ══ */
  @media (max-width: 860px) {
    .hero-grid {
      grid-template-columns: 1fr;
      gap: 24px;
    }
    .top-nav {
      flex-wrap: wrap;
      gap: 12px;
      border-radius: 24px;
      padding: 14px 18px;
    }
    .hero-headline {
      font-size: clamp(32px, 5.5vw, 42px);
    }
    .balance-card {
      padding: 20px 22px;
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
    .hero-card {
      padding: 20px 18px;
      border-radius: 22px;
      margin-bottom: 22px;
    }
    .hero-headline {
      font-size: 30px;
      letter-spacing: -0.5px;
    }
    .hero-subtitle {
      font-size: 14px;
    }
    .hero-actions {
      flex-direction: column;
      gap: 10px;
      width: 100%;
    }
    .btn-voyavera-primary,
    .btn-voyavera-glass {
      width: 100%;
      justify-content: center;
      padding: 13px 20px;
    }
    .weather-floating-card {
      padding: 14px 16px;
    }
    .balance-card {
      padding: 18px 16px;
      border-radius: 18px;
    }
    .bc-number {
      font-size: clamp(28px, 8vw, 38px);
    }
    .bc-footer {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
    .activity-section {
      padding: 0 4px;
    }
    .activity-header {
      flex-direction: column;
      align-items: stretch;
      gap: 14px;
      margin-bottom: 16px;
    }
    .activity-tools {
      width: 100%;
      display: flex;
      gap: 8px;
    }
    .search-input-wrap {
      flex: 1;
    }
    .entry-row {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;
      padding: 14px;
      border-radius: 18px;
    }
    .er-right {
      width: 100%;
      justify-content: space-between;
      border-top: 1px solid rgba(255, 255, 255, 0.08);
      padding-top: 10px;
    }
    .modal-backdrop {
      align-items: flex-end;
      padding: 0;
    }
    .modal-dialog {
      max-width: 100%;
      width: 100%;
      border-radius: 26px 26px 0 0;
      padding: 20px 18px calc(24px + env(safe-area-inset-bottom, 16px));
      max-height: 88vh;
      overflow-y: auto;
    }
    .quick-pill-grid {
      grid-template-columns: repeat(3, 1fr);
      gap: 8px;
    }
    .quick-pill {
      padding: 10px 4px;
      font-size: 13px;
    }
    .modal-footer {
      flex-direction: column-reverse;
      gap: 10px;
    }
    .modal-footer button {
      width: 100%;
      justify-content: center;
      padding: 13px;
    }
  }

  @media (max-width: 420px) {
    .nav-pill span {
      font-size: 12px;
    }
    .hero-headline {
      font-size: 26px;
    }
    .er-date-box {
      width: 44px;
      height: 44px;
    }
    .er-day {
      font-size: 16px;
    }
    .er-amount {
      font-size: 18px;
    }
  }
</style>
