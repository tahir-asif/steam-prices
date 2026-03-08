import { render, screen, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import App from './App'

// Create a mock fetch function that returns a promise which never resolves
// (so it doesn't trigger .then() unless we explicitly mock it)
const mockFetch = vi.fn(() => new Promise(() => {}))

// Replace global fetch with our mock
vi.stubGlobal('fetch', mockFetch)

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows checking connection initially', () => {
    render(<App />)
    // The fetch call is made, but it's a pending promise, so we still see the initial text
    expect(screen.getByText(/Checking connection/i)).toBeInTheDocument()
  })

  it('displays backend status when fetch succeeds', async () => {
    // Override the mock for this test
    mockFetch.mockResolvedValueOnce({
      json: async () => ({ status: 'ok' }),
    } as Response)

    render(<App />)

    await waitFor(() => {
      expect(screen.getByText(/Backend says: ok/i)).toBeInTheDocument()
    })
  })

  it('shows error when backend is unreachable', async () => {
    mockFetch.mockRejectedValueOnce(new Error('Network error'))

    render(<App />)

    await waitFor(() => {
      expect(screen.getByText(/Error: Could not connect/i)).toBeInTheDocument()
    })
  })
})
