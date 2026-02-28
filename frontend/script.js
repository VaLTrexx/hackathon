const DEFAULT_COUNTRIES = ["US", "IN", "CN", "IR", "JP"];
const BACKEND_URL = "http://localhost:8080/risk"; // change if needed

const grid = document.getElementById("countriesGrid");

function getRiskColor(score) {
    if (score < 30) return "green";
    if (score < 60) return "yellow";
    if (score < 80) return "orange";
    return "red";
}

function createCard(data) {
    const card = document.createElement("div");
    card.className = "card";
    card.style.borderColor = getRiskColor(data.risk_score);

    card.innerHTML = `
        <h2>${data.country}</h2>
        <div class="score">${data.risk_score}</div>
        <p>Volatility: ${data.volatility}</p>
        <p>Z-Score: ${data.z_score}</p>
    `;

    grid.appendChild(card);
}

async function fetchRisk(country) {
    try {
        const res = await fetch(`${BACKEND_URL}?country=${country}`);
        const data = await res.json();
        createCard(data);
    } catch (err) {
        console.error("Error fetching:", err);
    }
}

function searchCountry() {
    const input = document.getElementById("countryInput");
    const country = input.value.trim().toUpperCase();
    if (country) {
        fetchRisk(country);
        input.value = "";
    }
}

function loadDefaults() {
    DEFAULT_COUNTRIES.forEach(fetchRisk);
}

window.onload = loadDefaults;