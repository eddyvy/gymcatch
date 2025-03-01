import { Button } from '@mui/material'
import AppBar from '@mui/material/AppBar'
import Box from '@mui/material/Box'
import Toolbar from '@mui/material/Toolbar'
import Typography from '@mui/material/Typography'
import { FC } from 'react'

type Props = {
  setIsAuthenticated: (value: boolean) => void
}

export const Header: FC<Props> = ({ setIsAuthenticated }) => {
  const email = localStorage.getItem('email')

  const handleLogout = () => {
    localStorage.removeItem('sessionId')
    setIsAuthenticated(false)
  }

  return (
    <Box sx={{ flexGrow: 1 }} component="header">
      <AppBar position="static">
        {email && (
          <Toolbar>
            <Typography variant="caption" component="div" sx={{ flexGrow: 1 }}>
              {email}
            </Typography>
            <Button color="inherit" onClick={handleLogout}>
              Logout
            </Button>
          </Toolbar>
        )}
      </AppBar>
    </Box>
  )
}
