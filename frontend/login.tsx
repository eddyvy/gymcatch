import { useState, FC } from 'react'
import { TextField, Button, Typography, Box } from '@mui/material'
import { postAuth } from './features/client'

type Props = {
  setIsAuthenticated: (value: boolean) => void
  setIsLoading: (value: boolean) => void
  setError: (value: string | null) => void
}

export const Login: FC<Props> = ({
  setIsAuthenticated,
  setIsLoading,
  setError,
}) => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const handleLogin = () => {
    setIsLoading(true)
    postAuth(email, password)
      .then((sessionId) => {
        localStorage?.setItem('sessionId', sessionId)
        localStorage?.setItem('email', email)
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
      <Button type="submit" variant="contained" color="primary">
        Login
      </Button>
    </Box>
  )
}
