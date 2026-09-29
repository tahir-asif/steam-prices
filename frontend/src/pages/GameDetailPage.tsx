import { useParams, Link } from 'react-router-dom'
import { useState, useEffect } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  useActiveTooltipDataPoints,
  useActiveTooltipCoordinate,
  useCartesianScale,
  usePlotArea,
} from 'recharts'
import {
  getGamePriceHistory,
  type GamePriceData,
  type PricePoint,
} from '../services/api'
import styles from './GameDetailPage.module.css'
import LoadingSpinner from '../components/LoadingSpinner'

const DAY_MS = 24 * 60 * 60 * 1000

const toTime = (value: string) => new Date(value).getTime()

// Format cents to a currency string with a space between symbol and value
const formatPrice = (cents: number, currency: string) => {
  const code = currency || 'USD'
  const value = new Intl.NumberFormat(undefined, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(cents / 100)

  try {
    const symbol =
      new Intl.NumberFormat(undefined, { style: 'currency', currency: code })
        .formatToParts(0)
        .find((part) => part.type === 'currency')?.value ?? code
    return `${symbol} ${value}`
  } catch {
    return `${code} ${value}`
  }
}

// Helper to format date
const formatDate = (value: string | number) => {
  const date = new Date(value)
  return date.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

// Round a maximum value up to a visually clean axis top
const niceMax = (value: number) => {
  if (value <= 0) return 100
  const padded = value * 1.1
  const magnitude = 10 ** Math.floor(Math.log10(padded))
  const normalized = padded / magnitude
  const steps = [1, 1.5, 2, 2.5, 3, 4, 5, 6, 8, 10]
  const step = steps.find((candidate) => normalized <= candidate) ?? 10
  return step * magnitude
}

type RangeKey = '7D' | '30D' | '3M' | '6M' | 'YTD' | 'ALL'

const RANGES: RangeKey[] = ['7D', '30D', '3M', '6M', 'YTD', 'ALL']

const RANGE_DAYS: Record<string, number> = {
  '7D': 7,
  '30D': 30,
  '3M': 90,
  '6M': 182,
}

const rangeCutoff = (range: RangeKey, now: Date): number | null => {
  if (range === 'ALL') return null
  if (range === 'YTD') return new Date(now.getFullYear(), 0, 1).getTime()
  return now.getTime() - RANGE_DAYS[range] * DAY_MS
}

interface ChartPoint extends PricePoint {
  time: number
  formattedPrice: string
  formattedDate: string
  synthetic?: boolean
}

// Only real price changes get a dot; synthetic boundary points just extend the line
const renderDot = ({
  cx,
  cy,
  payload,
}: {
  cx?: number
  cy?: number
  payload?: ChartPoint
}) => {
  if (cx == null || cy == null || payload?.synthetic) {
    return <g />
  }
  return <circle className={styles.chartDot} cx={cx} cy={cy} r={4} />
}

// A hover label anchored to the active data point (not the mouse position)
function PointLabel() {
  const activePoints = useActiveTooltipDataPoints<ChartPoint>()
  const pointer = useActiveTooltipCoordinate()
  const plotArea = usePlotArea()
  const point = activePoints && activePoints.length > 0 ? activePoints[0] : undefined
  const pixel = useCartesianScale(
    point ? { x: point.time, y: point.price } : { x: 0, y: 0 },
  )

  if (!point || !pointer || !pixel) return null

  const PROXIMITY = 24
  if (Math.abs(pointer.y - pixel.y) > PROXIMITY) return null

  // Always place the label diagonally off the point, flipping sides near edges
  const EDGE_X = 80
  const EDGE_Y = 48
  const OFFSET_X = 16
  let anchor: 'start' | 'end' = 'start'
  let labelX = pixel.x + OFFSET_X
  if (plotArea) {
    const distLeft = pixel.x - plotArea.x
    const distRight = plotArea.x + plotArea.width - pixel.x
    if (distRight < EDGE_X && distRight <= distLeft) {
      anchor = 'end'
      labelX = pixel.x - OFFSET_X
    }
  }

  const placeBelow = !!plotArea && pixel.y - plotArea.y < EDGE_Y
  const dateY = placeBelow ? pixel.y + 40 : pixel.y - 32
  const priceY = placeBelow ? pixel.y + 24 : pixel.y - 16

  return (
    <g className={styles.pointLabel}>
      <text x={labelX} y={dateY} textAnchor={anchor} className={styles.pointLabelDate}>
        {point.formattedDate}
      </text>
      <text x={labelX} y={priceY} textAnchor={anchor} className={styles.pointLabelText}>
        {point.formattedPrice}
      </text>
    </g>
  )
}

function GameDetailPage() {
  const { appid } = useParams<{ appid: string }>()
  const [data, setData] = useState<GamePriceData | null>(null)
  const [range, setRange] = useState<RangeKey>('ALL')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [heroFailed, setHeroFailed] = useState(false)

  useEffect(() => {
    const fetchHistory = async () => {
      if (!appid) return

      try {
        setLoading(true)
        const result = await getGamePriceHistory(parseInt(appid, 10))
        setData(result)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load price history')
      } finally {
        setLoading(false)
      }
    }

    fetchHistory()
  }, [appid])

  if (loading) {
    return <LoadingSpinner message="Fetching price history..." />
  }

  if (error) {
    return <div className={styles.error}>Error: {error}</div>
  }

  const gameName = data?.name || `Game ${appid}`

  // Sort records oldest first, then collapse runs of the same price
  const rawHistory = [...(data?.history ?? [])].sort(
    (a, b) => toTime(a.recorded_at) - toTime(b.recorded_at),
  )
  const changePoints: PricePoint[] = []
  for (const point of rawHistory) {
    const previous = changePoints.at(-1)
    if (
      !previous ||
      previous.price !== point.price ||
      previous.currency !== point.currency
    ) {
      changePoints.push(point)
    }
  }

  // Base timeline: real price changes plus a synthetic point for today
  const today = new Date()
  const prepared: (PricePoint & { synthetic: boolean })[] = changePoints.map(
    (point) => ({ ...point, synthetic: false }),
  )
  const lastChange = prepared.at(-1)
  if (
    lastChange &&
    new Date(lastChange.recorded_at).toDateString() !== today.toDateString()
  ) {
    prepared.push({
      price: lastChange.price,
      currency: lastChange.currency,
      recorded_at: today.toISOString(),
      synthetic: true,
    })
  }

  const singlePoint = prepared.length === 1

  const currency = prepared.at(-1)?.currency || 'USD'
  const currentPrice = prepared.at(-1)?.price ?? 0
  const prices = prepared.map((point) => point.price)
  const maxPrice = prices.length > 0 ? Math.max(...prices) : 0
  const minPrice = prices.length > 0 ? Math.min(...prices) : 0
  const hasVariation = maxPrice !== minPrice
  const atRecordedLow = hasVariation && currentPrice === minPrice
  const atRecordedHigh = hasVariation && currentPrice === maxPrice
  const yMax = niceMax(maxPrice)

  // Trim to the selected range, adding a synthetic point at the range start
  const cutoff = rangeCutoff(range, today)
  let visible = prepared
  if (cutoff !== null) {
    const pointsInRange = prepared.filter((point) => toTime(point.recorded_at) >= cutoff)
    const cutoffDay = new Date(cutoff).toDateString()
    const hasCutoffPoint = pointsInRange.some(
      (point) => new Date(point.recorded_at).toDateString() === cutoffDay,
    )
    const priceAtCutoff = changePoints
      .filter((point) => toTime(point.recorded_at) <= cutoff)
      .at(-1)?.price

    visible = pointsInRange
    if (priceAtCutoff !== undefined && !hasCutoffPoint) {
      visible = [
        {
          price: priceAtCutoff,
          currency,
          recorded_at: new Date(cutoff).toISOString(),
          synthetic: true,
        },
        ...pointsInRange,
      ]
    }
  }

  const chartPoints: ChartPoint[] = visible.map((point) => ({
    ...point,
    time: toTime(point.recorded_at),
    formattedPrice: formatPrice(point.price, point.currency),
    formattedDate: formatDate(point.recorded_at),
  }))

  let xDomain: [number, number] = [0, 1]
  if (chartPoints.length === 1) {
    xDomain = [chartPoints[0].time - DAY_MS, chartPoints[0].time + DAY_MS]
  } else if (chartPoints.length > 1) {
    xDomain = [chartPoints[0].time, chartPoints[chartPoints.length - 1].time]
  }

  return (
    <div className={styles.page}>
      {appid && !heroFailed && (
        <div className={styles.hero}>
          <img
            src={`https://cdn.cloudflare.steamstatic.com/steam/apps/${appid}/library_hero.jpg`}
            alt=""
            className={styles.heroImage}
            onError={() => setHeroFailed(true)}
          />
        </div>
      )}

      <div className={styles.content}>
        <div className={styles.container}>
          <Link to="/" className={styles.backLink}>← Back to Search</Link>

          <div className={styles.header}>
            <h2 className={styles.title}>{gameName}</h2>
            <p className={styles.appId}>Steam App ID: {appid}</p>
          </div>

          {prepared.length === 0 ? (
            <p>No price history available yet. Check back later!</p>
          ) : (
            <>
              <div className={styles.stats}>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Current Price</span>
                  <span className={styles.statValue}>
                    {formatPrice(currentPrice, currency)}
                    {atRecordedLow && (
                      <span className={`${styles.badge} ${styles.badgeLow}`}>
                        Recorded low
                      </span>
                    )}
                    {atRecordedHigh && (
                      <span className={`${styles.badge} ${styles.badgeHigh}`}>
                        Recorded high
                      </span>
                    )}
                  </span>
                </div>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Highest Price</span>
                  <span className={styles.statValue}>
                    {formatPrice(maxPrice, currency)}
                  </span>
                </div>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Lowest Price</span>
                  <span className={styles.statValue}>
                    {formatPrice(minPrice, currency)}
                  </span>
                </div>
              </div>

              <div className={styles.chartContainer}>
                <ResponsiveContainer width="100%" height={400}>
                  <LineChart
                    data={chartPoints}
                    margin={{ top: 16, right: 16, bottom: 8, left: 4 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis
                      dataKey="time"
                      type="number"
                      scale="time"
                      domain={xDomain}
                      tickFormatter={(value) => formatDate(value)}
                      tick={{ fontSize: 12 }}
                      tickMargin={14}
                      minTickGap={32}
                    />
                    <YAxis
                      tickFormatter={(value) => formatPrice(value, currency)}
                      tick={{ fontSize: 12 }}
                      domain={[0, yMax]}
                      tickCount={5}
                      allowDecimals={false}
                      width={70}
                      tickMargin={8}
                    />
                    <Tooltip content={() => null} cursor={false} />
                    <Line
                      type="stepAfter"
                      dataKey="price"
                      stroke="currentColor"
                      strokeWidth={2}
                      dot={renderDot}
                      activeDot={{ r: 6 }}
                      isAnimationActive={false}
                    />
                    <PointLabel />
                  </LineChart>
                </ResponsiveContainer>
                {singlePoint && (
                  <p className={styles.note}>
                    We just started tracking this game. More data will appear over time.
                  </p>
                )}
              </div>

              <div className={styles.ranges}>
                {RANGES.map((option) => (
                  <button
                    key={option}
                    type="button"
                    onClick={() => setRange(option)}
                    className={`${styles.rangeButton} ${
                      range === option ? styles.rangeButtonActive : ''
                    }`}
                  >
                    {option}
                  </button>
                ))}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default GameDetailPage
