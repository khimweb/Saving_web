<script>
  import { onMount } from 'svelte';

  // Cambodian Gold Units
  const UNITS = [
    { key: 'domlerng', khmer: 'ដំឡើង', english: 'Domlerng', short: 'Dom', grams: 37.5, icon: '🏅', color: '#f59e0b', desc: '1 Dom = 10 Chi = 37.5g' },
    { key: 'chi',      khmer: 'ជី',    english: 'Chi',      short: 'Chi', grams: 3.75, icon: '💛', color: '#fbbf24', desc: '1 Chi = 10 Hun = 3.75g' },
    { key: 'hun',      khmer: 'ហ៊ុន',   english: 'Hun',      short: 'Hun', grams: 0.375, icon: '✨', color: '#fde047', desc: '1 Hun = 0.375g' },
    { key: 'g',        khmer: 'ក្រាម',  english: 'Gram',     short: 'g',   grams: 1.0,   icon: '⚖️', color: '#38bdf8', desc: 'Metric 1.0 Gram' },
    { key: 'baht',     khmer: 'បាត',   english: 'Baht',     short: 'Baht',grams: 15.244,icon: '🔶', color: '#fb923c', desc: 'Thai Baht = 15.244g' },
    { key: 'kg',       khmer: 'គីឡូ',  english: 'Kilogram', short: 'kg',  grams: 1000.0,icon: '👑', color: '#a855f7', desc: '1 Kg = 26.67 Dom' },
  ];

  let entries = [];
  let selectedUnit = UNITS[0]; // Default: Domlerng
  let qtyInput = '';
  let moneyInput = '';
  let noteInput = '';
  let datePicker = new Date().toISOString().slice(0, 10);
  let showModal = false;
  let loading = true;
  let timeStr = '';
  let dateStr = '';
  let searchQuery = '';
  let alert = { show: false, msg: '', type: 'gold', icon: '🥇' };
  let alertTimer = null;

  function triggerAlert(msg, type = 'gold', icon = '🥇') {
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

  async function fetchGoldEntries() {
    try {
      const res = await fetch('/api/gold');
      const data = await res.json();
      entries = Array.isArray(data) ? data : [];
    } catch (e) {
      entries = [];
    } finally {
      loading = false;
    }
  }

  // Reactive calculations
  $: totalGrams = entries.reduce((acc, e) => acc + Number(e.grams || 0), 0);
  $: totalMoney = entries.reduce((acc, e) => acc + Number(e.money || 0), 0);

  // Live conversions for header summary
  $: totalDomlerng = (totalGrams / 37.5).toFixed(2);
  $: totalChi = (totalGrams / 3.75).toFixed(2);
  $: totalHun = (totalGrams / 0.375).toFixed(1);

  // Modal live calculation
  $: inputGrams = qtyInput && parseFloat(qtyInput) > 0 ? (parseFloat(qtyInput) * selectedUnit.grams).toFixed(3) : '0.000';
  $: inputChi = qtyInput && parseFloat(qtyInput) > 0 ? (parseFloat(inputGrams) / 3.75).toFixed(2) : '0.00';
  $: inputDomlerng = qtyInput && parseFloat(qtyInput) > 0 ? (parseFloat(inputGrams) / 37.5).toFixed(3) : '0.000';

  async function handleSave() {
    const q = parseFloat(qtyInput);
    if (!datePicker || isNaN(q) || q <= 0) {
      triggerAlert('Please enter valid gold quantity', 'error', '⚠️');
      return;
    }

    const m = parseFloat(moneyInput) || 0;

    try {
      const res = await fetch('/api/gold', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          date: datePicker,
          unit: selectedUnit.key,
          qty: q,
          money: m,
          note: noteInput.trim()
        })
      });

      if (res.ok) {
        triggerAlert(`Added ${q} ${selectedUnit.short} (${inputGrams}g) to Vault!`, 'gold', '🥇');
        closeModal();
        await fetchGoldEntries();
      } else {
        triggerAlert('Failed to save gold entry', 'error', '✕');
      }
    } catch (e) {
      triggerAlert('Network error saving gold', 'error', '✕');
    }
  }

  async function deleteGold(id) {
    if (!confirm('Are you sure you want to delete this gold deposit?')) return;
    try {
      const res = await fetch(`/api/gold/${id}`, { method: 'DELETE' });
      if (res.ok) {
        triggerAlert('Gold deposit removed', 'error', '🗑️');
        await fetchGoldEntries();
      }
    } catch (e) {
      triggerAlert('Could not delete entry', 'error', '✕');
    }
  }

  function openAddModal() {
    qtyInput = '';
    moneyInput = '';
    noteInput = '';
    datePicker = new Date().toISOString().slice(0, 10);
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  function selectUnitQuick(u) {
    selectedUnit = u;
  }

  function getUnitMeta(key) {
    return UNITS.find(u => u.key === key.toLowerCase()) || UNITS[0];
  }

  $: filteredEntries = entries.filter(e => {
    if (!searchQuery.trim()) return true;
    const q = searchQuery.toLowerCase();
    return (e.date && e.date.toLowerCase().includes(q)) ||
           (e.unit && e.unit.toLowerCase().includes(q)) ||
           (e.note && e.note.toLowerCase().includes(q)) ||
           (e.qty && e.qty.toString().includes(q));
  });

  onMount(() => {
    updateCambodiaTime();
    const timer = setInterval(updateCambodiaTime, 1000);
    fetchGoldEntries();
    return () => clearInterval(timer);
  });
</script>

<!-- ══ AMBIENT BACKGROUND SCENE (Voyavera Alpine Sky + Golden Halo) ══ -->
<div class="ambient-scene">
  <div class="ambient-orb ambient-orb-1"></div>
  <div class="ambient-orb ambient-orb-gold"></div>
  <div class="ambient-orb ambient-orb-3"></div>
</div>

<!-- ══ DYNAMIC ISLAND ALERT ══ -->
{#if alert.show}
  <div class="dynamic-island-alert {alert.type}">
    <span class="dia-icon">{alert.icon}</span>
    <span class="dia-text">{alert.msg}</span>
    <div class="dia-glow-bar"></div>
  </div>
{/if}

<main class="app-container">
  <!-- ══ TOP NAVIGATION (Voyavera Pill Header) ══ -->
  <header class="top-nav glass-panel">
    <div class="brand">
      <div class="brand-logo gold-glow">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
        </svg>
      </div>
      <span class="brand-name">KHIM X KAIYA</span>
    </div>

    <!-- Nav Pills -->
    <nav class="nav-pills">
      <a href="index.html" class="nav-pill">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
          <circle cx="12" cy="12" r="10"/><path d="M16 8h-6a2 2 0 1 0 0 4h4a2 2 0 1 1 0 4H8"/><path d="M12 18V6"/>
        </svg>
        <span>Savings</span>
      </a>
      <a href="gold.html" class="nav-pill active gold-pill">
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

    <!-- Cambodia Live Badge -->
    <div class="live-clock-badge">
      <span class="live-dot gold"></span>
      <div class="clock-col">
        <span class="clock-digits">{timeStr || '00:00:00'}</span>
        <span class="clock-label">🇰🇭 Gold Market</span>
      </div>
    </div>
  </header>

  <!-- ══ GOLD HERO BANNER ══ -->
  <section class="gold-hero-card glass-panel">
    <div class="hero-grid">
      <!-- Left Content -->
      <div class="hero-left">
        <div class="eyebrow-tag gold">
          <span class="tag-sparkle">✦</span>
          <span>CAMBODIAN PRECIOUS METALS RESERVE</span>
        </div>
        <h1 class="hero-headline">
          The Gold <br />
          <span class="gold-gradient-text">Awaits You</span>
        </h1>
        <p class="hero-subtitle">
          Secure savings denominated in traditional Cambodian units: <strong>ដំឡើង (Domlerng)</strong>, <strong>ជី (Chi)</strong>, and <strong>ហ៊ុន (Hun)</strong>.
        </p>

        <div class="hero-actions">
          <button class="btn-gold-primary" on:click={openAddModal}>
            <span>Deposit Gold</span>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </button>
          <a href="analysis.html" class="btn-voyavera-glass">
            <span>View Performance</span>
            <span>📊</span>
          </a>
        </div>
      </div>

      <!-- Right Multi-Unit Gold Card -->
      <div class="hero-right">
        <!-- Floating Conversion Pill Widget -->
        <div class="gold-weight-card">
          <div class="gwc-top">
            <span class="gwc-label">TOTAL ACCUMULATED GOLD</span>
            <span class="gwc-badge">{entries.length} Deposits</span>
          </div>

          <div class="gwc-main-weight">
            <span class="gwc-number">{totalGrams.toFixed(2)}</span>
            <span class="gwc-unit">Grams</span>
          </div>

          <!-- Cambodian Equivalents Grid -->
          <div class="cambodia-units-pills">
            <div class="cu-pill">
              <span class="cu-kh">ដំឡើង</span>
              <span class="cu-val">{totalDomlerng} Dom</span>
            </div>
            <div class="cu-pill">
              <span class="cu-kh">ជី</span>
              <span class="cu-val">{totalChi} Chi</span>
            </div>
            <div class="cu-pill">
              <span class="cu-kh">ហ៊ុន</span>
              <span class="cu-val">{totalHun} Hun</span>
            </div>
          </div>

          <!-- Total Money Invested Footer -->
          <div class="gwc-footer">
            <div class="gwc-invested">
              <span class="gi-label">Capital Invested:</span>
              <span class="gi-val">${totalMoney.toLocaleString('en-US', { minimumFractionDigits: 2 })} USD</span>
            </div>
            <span class="gwc-add-link" on:click={openAddModal}>+ Add Entry</span>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- ══ INTERACTIVE UNIT GUIDE BAR ══ -->
  <section class="unit-guide-bar glass-panel">
    <div class="ug-header">
      <span class="ug-title">Cambodian Gold Conversion Standards</span>
      <span class="ug-sub">Official jewelry standard weights</span>
    </div>
    <div class="ug-grid">
      {#each UNITS as u}
        <button
          class="ug-item"
          class:active={selectedUnit.key === u.key}
          on:click={() => { selectedUnit = u; openAddModal(); }}
        >
          <div class="ug-icon-box" style="background: {u.color}22; border-color: {u.color}55;">
            <span>{u.icon}</span>
          </div>
          <div class="ug-info">
            <div class="ug-khmer">{u.khmer} ({u.english})</div>
            <div class="ug-calc">1 {u.short} = <strong>{u.grams}g</strong></div>
          </div>
        </button>
      {/each}
    </div>
  </section>

  <!-- ══ GOLD TRANSACTIONS LEDGER ══ -->
  <section class="activity-section">
    <div class="activity-header">
      <div>
        <h2 class="section-heading">Gold Vault Ledger</h2>
        <p class="section-sub">Individual bullion & jewelry holdings</p>
      </div>

      <div class="activity-tools">
        <div class="search-input-wrap">
          <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
            <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            type="text"
            placeholder="Search gold entries..."
            bind:value={searchQuery}
            class="search-input"
          />
          {#if searchQuery}
            <button class="clear-search" on:click={() => searchQuery = ''}>✕</button>
          {/if}
        </div>
        <button class="btn-add-circle gold" on:click={openAddModal} title="Deposit Gold">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Gold List -->
    {#if loading}
      <div class="skeleton-wrap">
        {#each [1, 2, 3] as _}
          <div class="skeleton-card"></div>
        {/each}
      </div>
    {:else if filteredEntries.length === 0}
      <div class="empty-glass-card glass-panel">
        <div class="empty-icon-halo gold">
          <svg width="44" height="44" viewBox="0 0 24 24" fill="none" stroke="#fbbf24" stroke-width="1.8">
            <path d="M6 3h12l4 6-10 12L2 9l4-6z"/><path d="M11 3 8 9l4 12 4-12-3-6"/>
          </svg>
        </div>
        <h3 class="empty-title">{searchQuery ? 'No matching gold holdings' : 'No gold deposited in vault'}</h3>
        <p class="empty-desc">
          {searchQuery ? 'Try another search keyword' : 'Deposit your first Domlerng, Chi, or Hun to start your gold reserve'}
        </p>
        <button class="btn-gold-primary sm" on:click={openAddModal}>
          <span>+ Deposit Gold</span>
        </button>
      </div>
    {:else}
      <div class="entries-grid">
        {#each filteredEntries as entry, index}
          {@const uMeta = getUnitMeta(entry.unit)}
          <div class="gold-entry-row glass-panel" style="animation-delay: {index * 35}ms">
            <!-- Left Info -->
            <div class="ger-left">
              <div class="ger-unit-box" style="background: {uMeta.color}15; border-color: {uMeta.color}45;">
                <span class="ger-icon">{uMeta.icon}</span>
                <span class="ger-khmer" style="color: {uMeta.color};">{uMeta.khmer}</span>
              </div>
              <div class="ger-details">
                <div class="ger-headline">
                  <span class="ger-qty">{entry.qty} {uMeta.english}</span>
                  <span class="ger-grams-tag">{Number(entry.grams).toFixed(3)} Grams</span>
                </div>
                <div class="ger-meta">
                  <span>📅 {entry.date}</span>
                  {#if entry.note}
                    <span class="ger-sep">•</span>
                    <span class="ger-note">{entry.note}</span>
                  {/if}
                </div>
              </div>
            </div>

            <!-- Right Money & Action -->
            <div class="ger-right">
              {#if entry.money && entry.money > 0}
                <div class="ger-money-wrap">
                  <span class="ger-money-label">PAID</span>
                  <span class="ger-money-val">${Number(entry.money).toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
                </div>
              {/if}
              <button class="action-btn delete" on:click={() => deleteGold(entry.id)} title="Delete">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
                  <polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </section>
</main>

<!-- ══ DEPOSIT GOLD MODAL SHEET ══ -->
{#if showModal}
  <div class="modal-backdrop" on:click={closeModal}>
    <div class="modal-dialog glass-panel gold-modal" on:click|stopPropagation>
      <div class="modal-top">
        <div class="modal-title-group">
          <div class="modal-tag gold">GOLD RESERVE ASSET</div>
          <h3 class="modal-heading">Deposit Gold to Vault</h3>
        </div>
        <button class="modal-close" on:click={closeModal}>✕</button>
      </div>

      <div class="modal-body">
        <!-- Unit Selection Pills -->
        <div class="form-field">
          <label class="ff-label">Select Cambodian Unit</label>
          <div class="unit-selector-grid">
            {#each UNITS as u}
              <button
                type="button"
                class="unit-pill-btn"
                class:selected={selectedUnit.key === u.key}
                on:click={() => selectUnitQuick(u)}
                style="--u-color: {u.color}"
              >
                <span class="upb-icon">{u.icon}</span>
                <span class="upb-name">{u.khmer}</span>
                <span class="upb-sub">{u.short} ({u.grams}g)</span>
              </button>
            {/each}
          </div>
        </div>

        <!-- Quantity & Date Row -->
        <div class="form-row-2">
          <div class="form-field">
            <label class="ff-label">Quantity ({selectedUnit.english})</label>
            <input
              type="number"
              class="ff-input gold-focus"
              placeholder="e.g. 1.5"
              step="0.001"
              min="0.001"
              bind:value={qtyInput}
              autofocus
            />
          </div>
          <div class="form-field">
            <label class="ff-label">Deposit Date</label>
            <input type="date" class="ff-input" bind:value={datePicker} />
          </div>
        </div>

        <!-- Live Instant Weight Calculator Preview -->
        <div class="live-calc-box">
          <div class="lcb-header">
            <span>⚡ LIVE CONVERSION</span>
            <span>1 {selectedUnit.short} = {selectedUnit.grams}g</span>
          </div>
          <div class="lcb-results">
            <div class="lcb-col">
              <span class="lcb-val">{inputGrams}g</span>
              <span class="lcb-lbl">Weight</span>
            </div>
            <div class="lcb-divider"></div>
            <div class="lcb-col">
              <span class="lcb-val">{inputChi} ជី</span>
              <span class="lcb-lbl">Chi</span>
            </div>
            <div class="lcb-divider"></div>
            <div class="lcb-col">
              <span class="lcb-val">{inputDomlerng} ដំ</span>
              <span class="lcb-lbl">Domlerng</span>
            </div>
          </div>
        </div>

        <!-- Money Paid (USD) -->
        <div class="form-field">
          <label class="ff-label">Total Money Paid (USD)</label>
          <div class="prefix-wrap">
            <span class="ff-prefix">$</span>
            <input
              type="number"
              class="ff-input prefix gold-focus"
              placeholder="0.00"
              step="0.01"
              min="0"
              bind:value={moneyInput}
            />
          </div>
        </div>

        <!-- Notes -->
        <div class="form-field">
          <label class="ff-label">Notes (Optional)</label>
          <input
            type="text"
            class="ff-input"
            placeholder="e.g. 24K Necklace, Phsar Thmei purchase..."
            bind:value={noteInput}
          />
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-cancel" on:click={closeModal}>Cancel</button>
        <button class="btn-gold-primary" on:click={handleSave}>
          <span>Save to Vault</span>
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
    background: linear-gradient(135deg, #f43f5e, #f59e0b);
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
    background: linear-gradient(135deg, #ffffff 40%, #fde68a 100%);
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
  .nav-pill.active.gold-pill {
    background: rgba(245, 158, 11, 0.2);
    color: #fde68a;
    border: 1px solid rgba(245, 158, 11, 0.45);
    box-shadow: 0 4px 18px rgba(245, 158, 11, 0.35);
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
  .live-dot.gold {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #fbbf24;
    box-shadow: 0 0 10px #fbbf24;
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

  /* ══ GOLD HERO BANNER ══ */
  .gold-hero-card {
    padding: 38px 42px;
    margin-bottom: 26px;
    border-radius: 36px;
    background: linear-gradient(135deg, rgba(30, 20, 10, 0.85) 0%, rgba(20, 14, 8, 0.8) 50%, rgba(10, 18, 38, 0.75) 100%);
    border: 1px solid rgba(245, 158, 11, 0.3);
    box-shadow: 0 25px 50px -15px rgba(0, 0, 0, 0.6), 0 0 40px rgba(245, 158, 11, 0.15);
  }
  .hero-grid {
    display: grid;
    grid-template-columns: 1.15fr 0.85fr;
    gap: 36px;
    align-items: center;
  }

  .eyebrow-tag.gold {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 1.5px;
    color: #fbbf24;
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
  .gold-gradient-text {
    background: linear-gradient(135deg, #fde68a 0%, #fbbf24 50%, #f59e0b 100%);
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
  .hero-subtitle strong {
    color: #fde68a;
  }

  .hero-actions {
    display: flex;
    gap: 14px;
    align-items: center;
  }

  .btn-gold-primary {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    background: linear-gradient(135deg, #fde68a 0%, #f59e0b 100%);
    color: #1c1917;
    padding: 13px 26px;
    border-radius: var(--radius-pill);
    font-size: 14.5px;
    font-weight: 800;
    border: none;
    cursor: pointer;
    box-shadow: 0 10px 25px -5px rgba(245, 158, 11, 0.45);
    transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  }
  .btn-gold-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 14px 30px -5px rgba(245, 158, 11, 0.65);
  }
  .btn-gold-primary:active {
    transform: translateY(0) scale(0.98);
  }
  .btn-gold-primary.sm {
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

  /* Right Gold Weight Card */
  .gold-weight-card {
    background: linear-gradient(135deg, rgba(60, 35, 15, 0.55) 0%, rgba(20, 15, 10, 0.65) 100%);
    border: 1.5px solid rgba(251, 191, 36, 0.4);
    border-radius: 28px;
    padding: 24px;
    box-shadow: 0 20px 45px -10px rgba(245, 158, 11, 0.25);
  }
  .gwc-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }
  .gwc-label {
    font-size: 11px;
    font-weight: 800;
    letter-spacing: 1.2px;
    color: #fbbf24;
  }
  .gwc-badge {
    background: rgba(251, 191, 36, 0.2);
    border: 1px solid rgba(251, 191, 36, 0.35);
    color: #fde68a;
    padding: 2px 10px;
    border-radius: var(--radius-pill);
    font-size: 11px;
    font-weight: 700;
  }

  .gwc-main-weight {
    display: flex;
    align-items: baseline;
    gap: 8px;
    margin-bottom: 16px;
  }
  .gwc-number {
    font-size: 44px;
    font-weight: 800;
    letter-spacing: -1.5px;
    color: #ffffff;
  }
  .gwc-unit {
    font-size: 18px;
    font-weight: 700;
    color: #fbbf24;
  }

  .cambodia-units-pills {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
    margin-bottom: 18px;
  }
  .cu-pill {
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 14px;
    padding: 10px 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 3px;
  }
  .cu-kh {
    font-size: 12px;
    font-weight: 700;
    color: #fde68a;
  }
  .cu-val {
    font-size: 13.5px;
    font-weight: 800;
    color: #ffffff;
  }

  .gwc-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 14px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
  }
  .gwc-invested {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .gi-label {
    font-size: 11px;
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .gi-val {
    font-size: 15px;
    font-weight: 800;
    color: #6ee7b7;
  }
  .gwc-add-link {
    font-size: 12.5px;
    font-weight: 700;
    color: #fbbf24;
    cursor: pointer;
    text-decoration: underline;
  }

  /* ══ UNIT GUIDE BAR ══ */
  .unit-guide-bar {
    padding: 20px 24px;
    border-radius: 24px;
    margin-bottom: 30px;
    background: rgba(15, 24, 48, 0.55);
  }
  .ug-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 14px;
  }
  .ug-title {
    font-size: 13px;
    font-weight: 800;
    letter-spacing: 0.8px;
    color: #fde68a;
    text-transform: uppercase;
  }
  .ug-sub {
    font-size: 12px;
    color: var(--text-tertiary);
  }
  .ug-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 10px;
  }
  .ug-item {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 16px;
    padding: 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    transition: all 0.2s;
    text-align: center;
  }
  .ug-item:hover {
    background: rgba(251, 191, 36, 0.12);
    border-color: rgba(251, 191, 36, 0.35);
    transform: translateY(-2px);
  }
  .ug-icon-box {
    width: 34px;
    height: 34px;
    border-radius: 10px;
    border: 1px solid;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
  }
  .ug-khmer {
    font-size: 12px;
    font-weight: 700;
    color: #fff;
  }
  .ug-calc {
    font-size: 11px;
    color: var(--text-secondary);
  }
  .ug-calc strong {
    color: #fbbf24;
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
    border-color: #fbbf24;
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

  .btn-add-circle.gold {
    width: 38px;
    height: 38px;
    border-radius: 50%;
    background: linear-gradient(135deg, #f59e0b, #fbbf24);
    border: none;
    color: #1c1917;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    box-shadow: 0 4px 14px rgba(245, 158, 11, 0.45);
    transition: all 0.2s;
  }
  .btn-add-circle.gold:hover {
    transform: scale(1.08);
  }

  /* Entries Grid */
  .entries-grid {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .gold-entry-row {
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

  .ger-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .ger-unit-box {
    width: 48px;
    height: 48px;
    border-radius: 14px;
    border: 1px solid;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    line-height: 1;
  }
  .ger-icon {
    font-size: 18px;
  }
  .ger-khmer {
    font-size: 10px;
    font-weight: 800;
    margin-top: 2px;
  }

  .ger-details {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }
  .ger-headline {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .ger-qty {
    font-size: 16px;
    font-weight: 800;
    color: #ffffff;
  }
  .ger-grams-tag {
    background: rgba(251, 191, 36, 0.15);
    border: 1px solid rgba(251, 191, 36, 0.3);
    color: #fde68a;
    padding: 2px 8px;
    border-radius: var(--radius-pill);
    font-size: 11px;
    font-weight: 700;
  }

  .ger-meta {
    font-size: 12px;
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .ger-sep {
    color: rgba(255, 255, 255, 0.2);
  }
  .ger-note {
    color: var(--text-secondary);
  }

  .ger-right {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .ger-money-wrap {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
  }
  .ger-money-label {
    font-size: 9.5px;
    font-weight: 700;
    letter-spacing: 0.8px;
    color: var(--text-tertiary);
  }
  .ger-money-val {
    font-size: 16px;
    font-weight: 800;
    color: #6ee7b7;
  }

  .action-btn.delete {
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
  .action-btn.delete:hover {
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
  .empty-icon-halo.gold {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: rgba(245, 158, 11, 0.1);
    border: 1px solid rgba(245, 158, 11, 0.35);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 18px;
    box-shadow: 0 0 30px rgba(245, 158, 11, 0.25);
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

  /* ══ MODAL SHEET ══ */
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(2, 6, 18, 0.75);
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
    max-width: 520px;
    background: rgba(18, 22, 38, 0.95);
    border: 1px solid rgba(245, 158, 11, 0.3);
    border-radius: 32px;
    padding: 30px;
    box-shadow: 0 25px 60px -15px rgba(0, 0, 0, 0.8), 0 0 40px rgba(245, 158, 11, 0.2);
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
    margin-bottom: 20px;
  }
  .modal-tag.gold {
    font-size: 10.5px;
    font-weight: 700;
    letter-spacing: 1.2px;
    color: #fbbf24;
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
    gap: 16px;
    margin-bottom: 24px;
  }

  /* Unit Selector Grid */
  .unit-selector-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }
  .unit-pill-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    padding: 10px 6px;
    background: rgba(255, 255, 255, 0.05);
    border: 1.5px solid rgba(255, 255, 255, 0.1);
    border-radius: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }
  .unit-pill-btn:hover {
    background: rgba(255, 255, 255, 0.09);
    border-color: rgba(255, 255, 255, 0.25);
  }
  .unit-pill-btn.selected {
    background: rgba(245, 158, 11, 0.18);
    border-color: #fbbf24;
    box-shadow: 0 4px 14px rgba(245, 158, 11, 0.3);
  }
  .upb-icon {
    font-size: 18px;
  }
  .upb-name {
    font-size: 12.5px;
    font-weight: 800;
    color: #fff;
  }
  .upb-sub {
    font-size: 10px;
    color: var(--text-tertiary);
  }

  .form-row-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
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
  .ff-input {
    width: 100%;
    background: rgba(255, 255, 255, 0.06);
    border: 1.5px solid rgba(255, 255, 255, 0.12);
    border-radius: 14px;
    padding: 12px 14px;
    font-size: 15px;
    color: #fff;
    font-family: inherit;
    outline: none;
    transition: all 0.2s;
  }
  .ff-input:focus.gold-focus {
    border-color: #fbbf24;
    background: rgba(245, 158, 11, 0.08);
    box-shadow: 0 0 0 4px rgba(245, 158, 11, 0.15);
  }

  .prefix-wrap {
    position: relative;
  }
  .ff-prefix {
    position: absolute;
    left: 14px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 17px;
    font-weight: 700;
    color: #fbbf24;
  }
  .ff-input.prefix {
    padding-left: 32px;
    font-weight: 700;
  }

  /* Live Calc Box */
  .live-calc-box {
    background: rgba(245, 158, 11, 0.08);
    border: 1px solid rgba(245, 158, 11, 0.25);
    border-radius: 16px;
    padding: 12px 16px;
  }
  .lcb-header {
    display: flex;
    justify-content: space-between;
    font-size: 10.5px;
    font-weight: 800;
    letter-spacing: 0.8px;
    color: #fbbf24;
    margin-bottom: 8px;
  }
  .lcb-results {
    display: flex;
    justify-content: space-around;
    align-items: center;
  }
  .lcb-col {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
  }
  .lcb-val {
    font-size: 16px;
    font-weight: 800;
    color: #fff;
  }
  .lcb-lbl {
    font-size: 10.5px;
    color: var(--text-tertiary);
  }
  .lcb-divider {
    width: 1px;
    height: 24px;
    background: rgba(255, 255, 255, 0.1);
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
    .ug-grid {
      grid-template-columns: repeat(3, 1fr);
      gap: 10px;
    }
    .hero-headline {
      font-size: clamp(32px, 5.5vw, 42px);
    }
    .gold-widget-card {
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
    .gold-hero-card {
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
    .btn-gold-primary,
    .btn-gold-glass {
      width: 100%;
      justify-content: center;
      padding: 13px 20px;
    }
    .gold-widget-card {
      padding: 18px 16px;
      border-radius: 18px;
    }
    .gwc-number {
      font-size: clamp(28px, 8vw, 38px);
    }
    .unit-breakdown-pills {
      gap: 6px;
    }
    .ub-pill {
      font-size: 11px;
      padding: 3px 8px;
    }
    .gwc-footer {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
    .unit-guide-card {
      padding: 18px 16px;
      border-radius: 20px;
      margin-bottom: 22px;
    }
    .ug-grid {
      grid-template-columns: repeat(2, 1fr);
      gap: 8px;
    }
    .ug-item {
      padding: 10px 8px;
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
    .gold-entry-row {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;
      padding: 14px;
      border-radius: 18px;
    }
    .ger-right {
      width: 100%;
      justify-content: space-between;
      border-top: 1px solid rgba(255, 255, 255, 0.08);
      padding-top: 10px;
    }
    .modal-backdrop {
      align-items: flex-end;
      padding: 0;
    }
    .gold-modal {
      max-width: 100%;
      width: 100%;
      border-radius: 26px 26px 0 0;
      padding: 20px 18px calc(24px + env(safe-area-inset-bottom, 16px));
      max-height: 88vh;
      overflow-y: auto;
    }
    .unit-selector-grid {
      grid-template-columns: repeat(3, 1fr);
      gap: 6px;
    }
    .form-row-2 {
      grid-template-columns: 1fr;
      gap: 12px;
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
    .ug-grid {
      grid-template-columns: 1fr;
    }
    .unit-selector-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }
</style>
