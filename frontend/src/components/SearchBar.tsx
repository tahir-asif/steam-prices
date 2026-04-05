import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { searchGames, type SearchResult } from '../services/api'
import styles from './SearchBar.module.css'

function SearchBar() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<SearchResult[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [showDropdown, setShowDropdown] = useState(false)
  const navigate = useNavigate()
  const debounceTimer = useRef<number | null>(null)

  // Effect to handle debounced search
  useEffect(() => {
    // Clear previous timer
    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current)
    }

    // Don't search if query is empty
    if (query.trim() === '') {
      setResults([])
      setShowDropdown(false)
      return
    }

    setIsLoading(true)

    // Set a new timer to call API after 300ms of no typing
    debounceTimer.current = setTimeout(async () => {
      try {
        const data = await searchGames(query)
        setResults(data)
        setShowDropdown(true)
      } catch (error) {
        console.error('Search error:', error)
        setResults([])
      } finally {
        setIsLoading(false)
      }
    }, 300)

    // Cleanup function to clear timer if component unmounts or query changes
    return () => {
      if (debounceTimer.current) {
        clearTimeout(debounceTimer.current)
      }
    }
  }, [query])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value)
  }

  const handleResultClick = (appid: number) => {
    setShowDropdown(false)
    setQuery('')
    navigate(`/game/${appid}`)
  }

  return (
    <div className={styles.searchContainer}>
      <input
        type="text"
        value={query}
        onChange={handleInputChange}
        placeholder="Search for a game..."
        className={styles.searchInput}
      />
      {isLoading && <div className={styles.loading}>Searching...</div>}
      {showDropdown && results.length > 0 && (
        <ul className={styles.dropdown}>
          {results.map((game) => (
            <li
              key={game.appid}
              onClick={() => handleResultClick(game.appid)}
              className={styles.dropdownItem}
            >
              <img src={game.icon} alt={game.name} width="32" height="32" />
              <span>{game.name}</span>
            </li>
          ))}
        </ul>
      )}
      {showDropdown && query && results.length === 0 && !isLoading && (
        <div className={styles.noResults}>No games found</div>
      )}
    </div>
  )
}

export default SearchBar
