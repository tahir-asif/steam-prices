// Use environment variable for API base URL, fallback to localhost for development
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

export interface SearchResult {
  appid: number
  name: string
  icon: string
}

export interface PricePoint {
  price: number
  currency: string
  recorded_at: string
}

export async function searchGames(query: string): Promise<SearchResult[]> {
  const response = await fetch(`${API_BASE}/search?q=${encodeURIComponent(query)}`)
  if (!response.ok) {
    throw new Error('Search failed')
  }
  return response.json()
}

export async function getPriceHistory(appId: number): Promise<PricePoint[]> {
  const response = await fetch(`${API_BASE}/games/${appId}/history`)
  if (!response.ok) {
    throw new Error('Failed to fetch price history')
  }
  return response.json()
}
