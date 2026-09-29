import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { searchGames, type SearchResult } from '../services/api'
import styles from './SearchBar.module.css'

function SearchBar() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<SearchResult[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [isSlow, setIsSlow] = useState(false)
  const [isOpen, setIsOpen] = useState(false)
  const navigate = useNavigate()
  const debounceTimer = useRef<number | null>(null)
  const slowTimer = useRef<number | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  // Effect to handle debounced search
  useEffect(() => {
    // Clear previous timers
    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current)
    }
    if (slowTimer.current) {
      clearTimeout(slowTimer.current)
    }

    // Don't search if query is empty
    if (query.trim() === '') {
      setResults([])
      setIsLoading(false)
      setIsSlow(false)
      return
    }

    setIsLoading(true)
    setIsSlow(false)

    // Show the cold-start hint if the request is still pending after 3s
    slowTimer.current = setTimeout(() => setIsSlow(true), 3000)

    // Call the API after 300ms of no typing
    debounceTimer.current = setTimeout(async () => {
      try {
        const data = await searchGames(query)
        setResults(data)
      } catch (error) {
        console.error('Search error:', error)
        setResults([])
      } finally {
        setIsLoading(false)
        setIsSlow(false)
        if (slowTimer.current) {
          clearTimeout(slowTimer.current)
        }
      }
    }, 300)

    // Cleanup function to clear timers if component unmounts or query changes
    return () => {
      if (debounceTimer.current) {
        clearTimeout(debounceTimer.current)
      }
      if (slowTimer.current) {
        clearTimeout(slowTimer.current)
      }
    }
  }, [query])

  // Close the results panel when clicking outside or pressing Escape
  useEffect(() => {
    const handleMouseDown = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false)
      }
    }

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleMouseDown)
    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.removeEventListener('mousedown', handleMouseDown)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setQuery(value)
    setIsOpen(value.trim() !== '')
  }

  const handleResultClick = (appid: number) => {
    setIsOpen(false)
    setQuery('')
    navigate(`/game/${appid}`)
  }

  return (
    <div className={styles.searchContainer} ref={containerRef}>
      <input
        type="text"
        value={query}
        onChange={handleInputChange}
        placeholder="Search for a game..."
        className={styles.searchInput}
      />
      {isOpen && (
        <div className={styles.panel}>
          {isLoading ? (
            <>
              <div className={styles.loadingRow}>
                <span className={styles.spinner} aria-hidden="true" />
                <span>Searching Steam...</span>
              </div>
              {isSlow && (
                <p className={styles.slowHint}>
                  Server is waking from cold start. This might take up to 1 minute.
                </p>
              )}
            </>
          ) : results.length > 0 ? (
            <ul className={styles.results}>
              {results.map((game) => (
                <li
                  key={game.appid}
                  onClick={() => handleResultClick(game.appid)}
                  className={styles.dropdownItem}
                >
                  <img src={game.icon} alt="" className={styles.dropdownThumb} />
                  <span>{game.name}</span>
                </li>
              ))}
            </ul>
          ) : (
            <div className={styles.noResults}>No games found</div>
          )}
        </div>
      )}
    </div>
  )
}

export default SearchBar
