import {
  Box,
  Typography,
  Button,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Tooltip,
} from '@mui/material'
import { FC, useEffect, useState, useMemo } from 'react'
import {
  getEvents,
  getInscribedEvents,
  getMegaEventsBooked,
  inscribeToEvent,
} from './features/client'
import { MegaEvent } from './features/events'
import { DateTime } from 'luxon'

type Props = {
  setIsLoading: (value: boolean) => void
  setError: (value: string | null) => void
}

export const Home: FC<Props> = ({ setError, setIsLoading }) => {
  const [events, setEvents] = useState<MegaEvent[]>([])
  const [inscribedEvents, setInscribedEvents] = useState<number[]>([])
  const [eventsBooked, setEventsBooked] = useState<number[]>([])
  const [selectedDate, setSelectedDate] = useState<string>('')

  useEffect(() => {
    setIsLoading(true)
    getEvents()
      .then((events) => {
        setEvents(events)
        return getInscribedEvents()
      })
      .then((eventIds) => {
        setInscribedEvents(eventIds)
      })
      .catch((err) => {
        setError(err.message)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }, [])

  const handleInscribe = (classId: number) => () => {
    setIsLoading(true)
    inscribeToEvent(classId)
      .then((success) => {
        if (success) {
          setInscribedEvents([...inscribedEvents, classId])
        } else {
          setError('Error inscribiendo a la clase')
        }
      })
      .catch((err) => {
        setError(err.message)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  const uniqueDates = useMemo(() => {
    return Array.from(new Set(events.map((event) => event.hour.split('T')[0])))
      .map((e) => DateTime.fromFormat(e, 'yyyy-MM-dd'))
      .filter(
        (d) =>
          d >=
          DateTime.now().set({ hour: 0, minute: 0, second: 0, millisecond: 0 })
      )
      .sort()
  }, [events])

  const filteredEvents = useMemo(() => {
    return events.filter(
      (event) =>
        DateTime.fromISO(event.hour).toFormat('yyyy-MM-dd') === selectedDate &&
        DateTime.fromISO(event.hour) >= DateTime.now()
    )
  }, [events, selectedDate])

  useEffect(() => {
    if (filteredEvents.length === 0) {
      setEventsBooked([])
      return
    }
    getMegaEventsBooked(filteredEvents.map((event) => event.session_id)).then(
      (events) => {
        setEventsBooked(events)
      }
    )
  }, [filteredEvents])

  const isBooked = (classId: number) =>
    inscribedEvents.includes(classId) || eventsBooked.includes(classId)

  return (
    <Box>
      <Typography variant="h4" gutterBottom align="center" sx={{ mt: 3 }}>
        GYM CATCH
      </Typography>
      <Box sx={{ p: 2 }}>
        <FormControl fullWidth margin="normal">
          <InputLabel id="date-select-label">Select Date</InputLabel>
          <Select
            labelId="date-select-label"
            value={selectedDate}
            onChange={(e) => setSelectedDate(e.target.value)}
            label="Select Date"
          >
            {uniqueDates.map((date, idx) => (
              <MenuItem key={idx} value={date.toFormat('yyyy-MM-dd')}>
                {date.setLocale('es').toFormat('ccc dd/MM')}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>
      <Box sx={{ p: 2 }}>
        {filteredEvents.map((event) => (
          <Tooltip
            key={event.session_id}
            title={
              isBooked(event.session_id)
                ? 'Ya estás inscrito en esta clase'
                : ''
            }
            arrow
          >
            <span>
              <Button
                variant="outlined"
                color="primary"
                onClick={handleInscribe(event.session_id)}
                fullWidth
                style={{ marginBottom: '10px' }}
                disabled={isBooked(event.session_id)}
              >
                <Box>
                  <Typography variant="button">
                    {event.activity_name}
                  </Typography>
                  <br />
                  <Typography variant="body2">
                    {DateTime.fromISO(event.start).toFormat('HH:mm')}
                  </Typography>
                </Box>
              </Button>
            </span>
          </Tooltip>
        ))}
      </Box>
    </Box>
  )
}
