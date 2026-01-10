import './Results.css'

function Results({ hostData }) {
  if (!hostData) return null

  const formatDuration = (ms) => {
    if (!ms) return 'N/A'
    return `${(ms / 1000).toFixed(2)} sec`
  }

  return (
    <div className="results-container">
      {hostData.engineVersion && (
        <div className="engine-version">
          Engine: {hostData.engineVersion}
        </div>
      )}

      {hostData.endpoints && hostData.endpoints.length > 0 && (
        <div className="endpoints-section">
          <table className="endpoints-table">
            <thead>
              <tr>
                <th>IP Address</th>
                <th>Duration</th>
                <th>Has Warnings</th>
                <th>Grade</th>
              </tr>
            </thead>
            <tbody>
              {hostData.endpoints.map((endpoint, index) => (
                <tr key={index} className={index % 2 === 0 ? 'row-even' : 'row-odd'}>
                  <td className="ip-cell">
                    {endpoint.ipAddress}
                  </td>
                  <td className="duration-cell">
                    {endpoint.duration ? formatDuration(endpoint.duration) : 'N/A'}
                  </td>
                  <td className="warnings-cell">
                    {endpoint.hasWarnings ? 'Yes' : 'No'}
                  </td>
                  <td className="grade-cell">
                    {endpoint.grade ? (
                      <span className="grade-text">
                        {endpoint.grade}
                      </span>
                    ) : 'N/A'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

export default Results
