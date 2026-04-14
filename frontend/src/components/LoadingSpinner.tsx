import { useState, useEffect } from 'react'
import styles from './LoadingSpinner.module.css'

interface LoadingSpinnerProps {
    message?: string
    showTimer?: boolean
}

function LoadingSpinner({ message = 'Loading...', showTimer = true }: LoadingSpinnerProps) {
    const [elapsed, setElapsed] = useState(0)

    useEffect(() => {
        if (!showTimer) return

        const start = Date.now()
        const interval = setInterval(() => {
            setElapsed(Math.floor((Date.now() - start) / 1000))
        }, 1000)

        return () => clearInterval(interval)
    }, [showTimer])

    return (
        <div className={styles.container}>
            <div className={styles.spinner}></div>
            <p className={styles.message}>{message}</p>
            {showTimer && elapsed > 0 && (
                <p className={styles.timer}>Elapsed: {elapsed} second{elapsed !== 1 ? 's' : ''}</p>
            )}
            <p className={styles.hint}>
                ⏳ Render free tier – first request may take up to 30 seconds.
            </p>
        </div>
    )
}

export default LoadingSpinner
