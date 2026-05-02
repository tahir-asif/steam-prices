import styles from './LoadingSpinner.module.css'

interface LoadingSpinnerProps {
    message?: string
}

function LoadingSpinner({ message = 'Loading...' }: LoadingSpinnerProps) {
    return (
        <div className={styles.container}>
            <div className={styles.spinner}></div>
            <p className={styles.message}>{message}</p>
        </div>
    )
}

export default LoadingSpinner
