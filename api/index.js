// api/index.js - Vercel Serverless API Handler
import fs from 'fs';
import path from 'path';

const TMP_FILE = path.join('/tmp', 'khim_kaiya_savings.json');

const UNIT_TO_GRAMS = {
  g: 1.0,
  chi: 3.75,
  domlerng: 37.5,
  hun: 0.375,
  baht: 15.244,
  kg: 1000.0,
};

// Default initial state
const defaultState = {
  entries: [
    { date: "2026-10-23", amount: 285, note: "Initial deposit" }
  ],
  gold: [
    { id: 1, date: "2026-10-23", unit: "g", qty: 2, grams: 2, money: 1015, note: "2g Pure Gold" },
    { id: 2, date: "2026-10-23", unit: "hun", qty: 2.6, grams: 0.975, money: 269, note: "Cambodian Hun" },
    { id: 3, date: "2026-10-23", unit: "hun", qty: 5.8, grams: 2.175, money: 350, note: "Gold deposit" }
  ],
  nextGoldId: 4
};

function readData() {
  try {
    if (fs.existsSync(TMP_FILE)) {
      const raw = fs.readFileSync(TMP_FILE, 'utf-8');
      return JSON.parse(raw);
    }
  } catch (e) {
    console.error('Error reading tmp store:', e);
  }
  return { ...defaultState };
}

function writeData(data) {
  try {
    fs.writeFileSync(TMP_FILE, JSON.stringify(data, null, 2), 'utf-8');
  } catch (e) {
    console.error('Error writing tmp store:', e);
  }
}

function wmoDescription(code) {
  if (code === 0) return "Clear sky";
  if (code <= 2) return "Partly cloudy";
  if (code === 3) return "Overcast";
  if (code <= 49) return "Fog";
  if (code <= 59) return "Drizzle";
  if (code <= 69) return "Rain";
  if (code <= 79) return "Snow";
  if (code <= 82) return "Rain showers";
  if (code <= 84) return "Snow showers";
  if (code <= 99) return "Thunderstorm";
  return "Unknown";
}

function wmoIcon(code) {
  if (code === 0) return "☀️";
  if (code <= 2) return "⛅";
  if (code === 3) return "☁️";
  if (code <= 49) return "🌫️";
  if (code <= 69) return "🌧️";
  if (code <= 79) return "❄️";
  if (code <= 84) return "🌦️";
  if (code <= 99) return "⛈️";
  return "🌡️";
}

async function fetchWeather() {
  try {
    const url = "https://api.open-meteo.com/v1/forecast?latitude=11.5564&longitude=104.9282&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m";
    const res = await fetch(url, { headers: { 'User-Agent': 'KhimKaiyaApp/1.0' } });
    const data = await res.json();
    if (data && data.current) {
      return {
        temp: data.current.temperature_2m,
        description: wmoDescription(data.current.weather_code),
        humidity: data.current.relative_humidity_2m,
        wind_speed: data.current.wind_speed_10m,
        icon: wmoIcon(data.current.weather_code)
      };
    }
  } catch (e) {
    console.error('Weather fetch error:', e);
  }
  return {
    temp: 28,
    description: "Partly cloudy",
    humidity: 75,
    wind_speed: 10,
    icon: "⛅"
  };
}

export default async function handler(req, res) {
  // CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
  res.setHeader('Content-Type', 'application/json');

  if (req.method === 'OPTIONS') {
    return res.status(204).end();
  }

  // Parse path
  const urlObj = new URL(req.url, `http://${req.headers.host || 'localhost'}`);
  const pathname = urlObj.pathname.replace(/\/$/, '');
  const data = readData();

  // 1. /api/total
  if (pathname === '/api/total' && req.method === 'GET') {
    const total = data.entries.reduce((sum, item) => sum + Number(item.amount || 0), 0);
    return res.status(200).json({ total });
  }

  // 2. /api/entries
  if (pathname === '/api/entries' && req.method === 'GET') {
    const sorted = [...data.entries].sort((a, b) => b.date.localeCompare(a.date));
    return res.status(200).json(sorted);
  }

  // 3. /api/entry (POST)
  if (pathname === '/api/entry' && req.method === 'POST') {
    let body = req.body;
    if (typeof body === 'string') {
      try { body = JSON.parse(body); } catch (_) {}
    }
    if (!body || !body.date) {
      return res.status(400).json({ error: 'bad request' });
    }
    const existingIndex = data.entries.findIndex(e => e.date === body.date);
    const newEntry = {
      date: body.date,
      amount: parseFloat(body.amount) || 0,
      note: (body.note || '').trim()
    };
    if (existingIndex >= 0) {
      data.entries[existingIndex] = newEntry;
    } else {
      data.entries.push(newEntry);
    }
    writeData(data);
    return res.status(200).json({ status: 'ok' });
  }

  // 4. /api/entry/:date (DELETE)
  if (pathname.startsWith('/api/entry/') && req.method === 'DELETE') {
    const date = pathname.replace('/api/entry/', '');
    data.entries = data.entries.filter(e => e.date !== date);
    writeData(data);
    return res.status(200).json({ status: 'ok' });
  }

  // 5. /api/gold (GET / POST)
  if (pathname === '/api/gold') {
    if (req.method === 'GET') {
      const sorted = [...data.gold].sort((a, b) => b.date.localeCompare(a.date) || b.id - a.id);
      return res.status(200).json(sorted);
    }
    if (req.method === 'POST') {
      let body = req.body;
      if (typeof body === 'string') {
        try { body = JSON.parse(body); } catch (_) {}
      }
      if (!body || !body.date || !body.unit || !body.qty) {
        return res.status(400).json({ error: 'bad request' });
      }
      const unit = body.unit;
      const factor = UNIT_TO_GRAMS[unit] || 1.0;
      const qty = parseFloat(body.qty) || 0;
      const grams = qty * factor;
      const newGold = {
        id: data.nextGoldId++,
        date: body.date,
        unit: unit,
        qty: qty,
        grams: grams,
        money: parseFloat(body.money) || 0,
        note: (body.note || '').trim()
      };
      data.gold.push(newGold);
      writeData(data);
      return res.status(200).json({ status: 'ok' });
    }
  }

  // 6. /api/gold/:id (DELETE)
  if (pathname.startsWith('/api/gold/') && req.method === 'DELETE') {
    const id = parseInt(pathname.replace('/api/gold/', ''), 10);
    data.gold = data.gold.filter(g => g.id !== id);
    writeData(data);
    return res.status(200).json({ status: 'ok' });
  }

  // 7. /api/weather (GET)
  if (pathname === '/api/weather' && req.method === 'GET') {
    const weather = await fetchWeather();
    return res.status(200).json(weather);
  }

  return res.status(404).json({ error: 'not found' });
}
