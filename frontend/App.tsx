import { useState } from 'react'
import { Container } from '@mui/material'
import { Header } from './components/header'
import { Login } from './login'
import { Home } from './home'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  return (
    <Container sx={{ paddingLeft: 0, paddingRight: 0 }}>
      <Header setIsAuthenticated={setIsAuthenticated} />
      {isAuthenticated ? (
        <Login setIsAuthenticated={setIsAuthenticated} />
      ) : (
        <Home />
      )}
    </Container>
  )
}

export default App
