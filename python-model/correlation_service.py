from fastapi import FastAPI
from pydantic import BaseModel, Field
from typing import List, Optional
import numpy as np
from scipy.stats import pearsonr
import datetime

app = FastAPI(title="Country Risk Score API", version="1.0")

# -----------------------------
# Data Models
# -----------------------------

class Event(BaseModel):
    tone: float = Field(..., description="AvgTone from GDELT")
    mentions: int = Field(..., ge=1, description="Number of mentions")
    severity_weight: float = Field(1.0, ge=0.0, description="Optional GoldsteinScale weight")
    date: Optional[datetime.date] = None  # Optional event date for recency decay

class RiskRequest(BaseModel):
    events: List[Event]
    prices: List[float] = Field(..., min_items=2, description="Market close prices, chronological order")

# -----------------------------
# Core Calculations
# -----------------------------

def compute_event_score(events: List[Event], decay_days: int = 7) -> float:
    """Compute aggregated event score with optional recency decay"""
    total_score = 0.0
    today = datetime.date.today()
    for e in events:
        decay = 1.0
        if e.date:
            days_ago = (today - e.date).days
            decay = np.exp(-days_ago / decay_days)
        score = abs(e.tone) * np.log(e.mentions + 1) * e.severity_weight * decay
        total_score += score
    return total_score

def compute_returns(prices: List[float]) -> np.ndarray:
    prices = np.array(prices)
    returns = np.diff(prices) / prices[:-1]
    return returns

def compute_z_score(returns: np.ndarray):
    mean = np.mean(returns)
    std = np.std(returns)
    if std == 0:
        return 0.0
    latest_return = returns[-1]
    z_score = (latest_return - mean) / std
    return z_score

def compute_volatility(returns: np.ndarray):
    return np.std(returns)

def compute_correlation(event_score: float, returns: np.ndarray):
    if len(returns) < 2:
        return 0.0
    event_series = np.full(len(returns), event_score)
    corr, _ = pearsonr(event_series, returns)
    if np.isnan(corr):
        return 0.0
    return corr

def normalize(value, min_val, max_val):
    return (value - min_val) / (max_val - min_val + 1e-8)

def compute_risk_score(event_score, volatility, z_score, correlation,
                       weights=None):
    """
    Compute final risk score with optional dynamic weights.
    Inputs can be raw; normalization included here for robustness.
    """
    if weights is None:
        weights = {"event": 0.35, "volatility": 0.35, "z_score": 0.2, "correlation": 0.1}

    # Normalize metrics to 0-1 for stable weighting
    norm_event = np.tanh(event_score / 10)           # Event score approx capped
    norm_vol = np.tanh(volatility * 100)            # Volatility scaled
    norm_z = np.tanh(abs(z_score))                  # Z-score normalized
    norm_corr = np.tanh(abs(correlation))           # Correlation normalized

    score = (
        weights["event"] * norm_event * 100 +
        weights["volatility"] * norm_vol * 100 +
        weights["z_score"] * norm_z * 100 +
        weights["correlation"] * norm_corr * 100
    )

    return min(score, 100)

# -----------------------------
# API Endpoint
# -----------------------------

@app.post("/compute_risk")
def compute_risk(request: RiskRequest):
    if not request.events:
        return {"error": "No events provided"}
    if len(request.prices) < 2:
        return {"error": "Not enough price data to compute returns"}

    # Compute components
    event_score = compute_event_score(request.events)
    returns = compute_returns(request.prices)
    z_score = compute_z_score(returns)
    volatility = compute_volatility(returns)
    correlation = compute_correlation(event_score, returns)
    risk_score = compute_risk_score(event_score, volatility, z_score, correlation)

    # Optional: per-event contributions
    event_contributions = [
        round(abs(e.tone) * np.log(e.mentions + 1) * e.severity_weight, 4)
        for e in request.events
    ]

    return {
        "event_score": round(event_score, 4),
        "z_score": round(z_score, 4),
        "volatility": round(volatility, 6),
        "correlation": round(correlation, 4),
        "risk_score": round(risk_score, 4),
        "event_contributions": event_contributions
    }