import { useState, useEffect, FC } from 'react'
import {
  TextField,
  Button,
  Typography,
  Box,
  Backdrop,
  CircularProgress,
} from '@mui/material'
import { getCheckSession, postAuth } from './features/client'

type Props = {
  setIsAuthenticated: (value: boolean) => void
}

export const Login: FC<Props> = ({ setIsAuthenticated }) => {
  const [isLoading, setIsLoading] = useState(false)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    const sessionId = localStorage.getItem('sessionId')
    if (sessionId) {
      setIsLoading(true)
      getCheckSession(sessionId)
        .then((success) => {
          if (success) {
            setIsAuthenticated(true)
          } else {
            localStorage.removeItem('sessionId')
            localStorage.removeItem('email')
          }
        })
        .catch(() => {
          localStorage.removeItem('sessionId')
          localStorage.removeItem('email')
        })
        .finally(() => setIsLoading(false))
    }
  }, [])

  const handleLogin = () => {
    setIsLoading(true)
    postAuth(email, password)
      .then((sessionId) => {
        localStorage.setItem('sessionId', sessionId)
        localStorage.setItem('email', email)
        setIsAuthenticated(true)
      })
      .catch(() => {
        setError('Invalid email or password')
      })
      .finally(() => setIsLoading(false))
  }

  return (
    <Box
      component="form"
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      sx={{
        marginTop: '20vh',
      }}
      onSubmit={(e) => {
        e.preventDefault()
        handleLogin()
      }}
    >
      <Typography variant="h4" gutterBottom>
        GYM CATCH
      </Typography>
      <TextField
        label="Email"
        name="email"
        variant="outlined"
        margin="normal"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <TextField
        label="Password"
        name="password"
        type="password"
        variant="outlined"
        margin="normal"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      {error && <Typography color="error">{error}</Typography>}
      <Button type="submit" variant="contained" color="primary">
        Login
      </Button>
      <Backdrop
        sx={(theme) => ({ color: '#fff', zIndex: theme.zIndex.drawer + 1 })}
        open={isLoading}
      >
        <CircularProgress color="inherit" />
      </Backdrop>
    </Box>
  )
}
