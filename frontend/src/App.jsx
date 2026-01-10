import { useState } from 'react'
import './App.css'
import Results from './Results'

function App() {
  const [domain, setDomain] = useState('')
  const [useCache, setUseCache] = useState(false)
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    if (!domain.trim()) {
      setError('Please enter a domain')
      return
    }

    setLoading(true)
    setError(null)
    setResult(null)

    try {
      const url = new URL('http://localhost:8080/analyze')
      url.searchParams.set('host', domain.trim())
      if (useCache) {
        url.searchParams.set('fromCache', 'true')
      }
      
      const response = await fetch(url.toString())
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(errorText || `Error: ${response.status}`)
      }

      const data = await response.json()
      setResult(data)
    } catch (err) {
      setError(err.message || 'Error analyzing domain')
    } finally {
      setLoading(false)
    }
  }

  const handleNewScan = () => {
    setResult(null)
    setError(null)
    setDomain('')
    setUseCache(false)
  }

  return (
    <div className={`app ${result ? 'has-results' : ''}`}>
      <div className="container">
        {!result && (
          <>
            <h1>SSL Labs Checker</h1>
            <form onSubmit={handleSubmit} className="form">
              <input
                type="text"
                placeholder="Enter domain (e.g.: example.com)"
                value={domain}
                onChange={(e) => setDomain(e.target.value)}
                className="domain-input"
                disabled={loading}
              />
              <button type="submit" className="submit-button" disabled={loading}>
                {loading ? 'Analyzing...' : 'Analyze'}
              </button>
            </form>
            <div className="cache-checkbox">
              <label>
                <input
                  type="checkbox"
                  checked={useCache}
                  onChange={(e) => setUseCache(e.target.checked)}
                  disabled={loading}
                />
                <span>Use SSL Labs server cache</span>
              </label>
            </div>

            {error && (
              <div className="error-message">
                {error}
              </div>
            )}
          </>
        )}

        {result && (
          <>
            <h1>SSL Labs Checker</h1>
            <div className="results-header">
              <button onClick={handleNewScan} className="new-scan-button">
                ←
              </button>
            </div>
            <Results hostData={result} />
          </>
        )}
      </div>
    </div>
  )
}

export default App
