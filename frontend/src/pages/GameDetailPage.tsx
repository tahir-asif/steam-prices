import { useParams, Link } from 'react-router-dom'
import { useState, useEffect, useRef } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import { getPriceHistory, type PricePoint } from '../services/api'
import styles from './GameDetailPage.module.css'
import LoadingSpinner from '../components/LoadingSpinner'

// Helper to format price from cents to dollars
const formatPrice = (cents: number, currency: string) => {
  return new Intl.NumberFormat(undefined, {
    style: 'currency',
    currency: currency || 'USD',
  }).format(cents / 100)
}

// Helper to format date
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function GameDetailPage() {
  const { appid } = useParams<{ appid: string }>()
  const [history, setHistory] = useState<PricePoint[]>([])
  const gameName = `Game ${appid}`
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showSlowLoading, setShowSlowLoading] = useState(false)
  const slowTimer = useRef<number | null>(null)

  useEffect(() => {
    const fetchHistory = async () => {
      if (!appid) return

      try {
        setLoading(true)
        setShowSlowLoading(false)

        slowTimer.current = setTimeout(() => {
          setShowSlowLoading(true)
        }, 2000)

        const data = await getPriceHistory(parseInt(appid, 10))
        setHistory(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load price history')
      } finally {
        setLoading(false)
        setShowSlowLoading(false)
        if (slowTimer.current) clearTimeout(slowTimer.current)
      }
    }

    fetchHistory()
  }, [appid])

  if (loading) {
    return (
      <LoadingSpinner
        message="Fetching price history..."
        showTimer={showSlowLoading}
      />
    )
  }

  if (error) {
    return <div className={styles.error}>Error: {error}</div>
  }

  // Prepare data for chart (add formatted values for tooltip)
  const chartData = history.map((point) => ({
    ...point,
    formattedPrice: formatPrice(point.price, point.currency),
    formattedDate: formatDate(point.recorded_at),
  }))

  const CustomTooltip = ({ active, payload }: any) => {
    if (!active || !payload || payload.length === 0) return null

    const data = payload[0].payload as PricePoint & { formattedPrice: string; formattedDate: string }

    return (
      <div style={{
        backgroundColor: 'white',
        padding: '10px',
        border: '1px solid #ccc',
        borderRadius: '4px',
      }}>
        <p style={{ margin: 0 }}>{data.formattedDate}</p>
        <p style={{ margin: 0, fontWeight: 'bold' }}>{data.formattedPrice}</p>
      </div>
    )
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Link to="/" className={styles.backLink}>← Back to Search</Link>
        <h2>{gameName}</h2>
        <p className={styles.appId}>Steam App ID: {appid}</p>
      </div>

      {history.length === 0 ? (
        <p>No price history available yet. Check back later!</p>
      ) : (
        <div className={styles.chartContainer}>
          <ResponsiveContainer width="100%" height={400}>
            <LineChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="recorded_at"
                tickFormatter={(value) => formatDate(value)}
                tick={{ fontSize: 12 }}
              />
              <YAxis
                tickFormatter={(value) => formatPrice(value, history[0]?.currency || 'USD')}
                tick={{ fontSize: 12 }}
                domain={['auto', 'auto']}
              />
              <Tooltip content={<CustomTooltip />} />
              <Line
                type="stepAfter"
                dataKey="price"
                stroke="#8884d8"
                strokeWidth={2}
                dot={{ r: 4 }}
                activeDot={{ r: 6 }}
              />
            </LineChart>
          </ResponsiveContainer>
          {history.length === 1 && (
            <p className={styles.note}>
              We just started tracking this game. More data will appear over time.
            </p>
          )}
        </div>
      )}
    </div>
  )
}

export default GameDetailPage
