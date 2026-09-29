import { useState, useEffect } from 'react'
import styles from './LoadingSpinner.module.css'

interface LoadingSpinnerProps {
  message?: string
}

const SLOW_THRESHOLD_MS = 3000

function LoadingSpinner({ message = 'Loading...' }: LoadingSpinnerProps) {
  const [isSlow, setIsSlow] = useState(false)

  useEffect(() => {
    const timer = setTimeout(() => setIsSlow(true), SLOW_THRESHOLD_MS)
    return () => clearTimeout(timer)
  }, [])

  return (
    <div className={styles.container}>
      <div className={styles.spinner}></div>
      <p className={styles.message}>{message}</p>
      {isSlow && (
        <p className={styles.hint}>
          Server is waking from cold start. This might take up to 1 minute.
        </p>
      )}
    </div>
  )
}

export default LoadingSpinner
