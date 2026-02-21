import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [backendStatus, setBackendStatus] = useState('Checking connection...')

  useEffect(() => {
    fetch('http://localhost:3000/api/health')
      .then(res => res.json())
      .then(data => setBackendStatus(`Backend says: ${data.status}`))
      .catch(() => setBackendStatus('Error: Could not connect to backend. Is it running?'))
  }, [])

  return (
    <>
      <h1>Steam Price Tracker</h1>
      <p>{backendStatus}</p>
    </>
  )
}

export default App
