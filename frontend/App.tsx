import { useEffect, useState } from 'react'
import {
  Backdrop,
  CircularProgress,
  Container,
  Paper,
  Typography,
} from '@mui/material'
import { Header } from './components/header'
import { Login } from './login'
import { Home } from './home'
import { getCheckSession } from './features/client'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const sessionId = localStorage?.getItem('sessionId')
    if (sessionId) {
      setIsLoading(true)
      getCheckSession(sessionId)
        .then((success) => {
          if (success) {
            setIsAuthenticated(true)
          } else {
            localStorage?.removeItem('sessionId')
            localStorage?.removeItem('email')
          }
        })
        .catch(() => {
          localStorage?.removeItem('sessionId')
          localStorage?.removeItem('email')
        })
        .finally(() => setIsLoading(false))
    }
  }, [])

  return (
    <Container
      sx={{
        paddingLeft: '0px !important',
        paddingRight: '0px !important',
      }}
    >
      <Header setIsAuthenticated={setIsAuthenticated} />
      {isAuthenticated ? (
        <Home setError={setError} setIsLoading={setIsLoading} />
      ) : (
        <Login
          setIsAuthenticated={setIsAuthenticated}
          setError={setError}
          setIsLoading={setIsLoading}
        />
      )}
      {error && (
        <Backdrop
          sx={(theme) => ({ color: '#fff', zIndex: theme.zIndex.drawer + 1 })}
          open={!!error}
          onClick={() => setError(null)}
        >
          <Paper sx={{ padding: 2 }}>
            <Typography color="error">{error}</Typography>
          </Paper>
        </Backdrop>
      )}
      {isLoading && (
        <Backdrop
          sx={(theme) => ({ color: '#fff', zIndex: theme.zIndex.drawer + 1 })}
          open={isLoading}
        >
          <CircularProgress color="inherit" />
        </Backdrop>
      )}
    </Container>
  )
}

export default App
