import { Box, Typography } from '@mui/material'
import { FC, useEffect, useState } from 'react'
import { getEvents } from './features/client'
import { MegaEvent } from './features/events'

type Props = {
  setIsLoading: (value: boolean) => void
  setError: (value: string | null) => void
}

export const Home: FC<Props> = ({ setError, setIsLoading }) => {
  const [events, setEvents] = useState<MegaEvent[]>([])

  useEffect(() => {
    getEvents()
      .then((events) => {
        setEvents(events)
        setIsLoading(false)
      })
      .catch((err) => {
        setError(err.message)
        setIsLoading(false)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }, [])

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Clases Mega
      </Typography>
      {events.map((event) => (
        <Typography key={event.id}>{event.title}</Typography>
      ))}
    </Box>
  )
}
