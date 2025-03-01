import axios from 'axios'
import { MegaEvent, EventResponse } from './events'

const BACKEND_HOST = import.meta.env.VITE_BACKEND_HOST || ''

export async function getCheckSession(sessionId: string): Promise<boolean> {
  const res = await axios.get<{ success: boolean }>(
    BACKEND_HOST + '/api/check_session/' + sessionId
  )
  return res.data.success
}

export async function postAuth(
  email: string,
  password: string
): Promise<string> {
  const res = await axios.post<{ sessionID: string }>(
    BACKEND_HOST + '/api/auth',
    {
      email,
      password,
    }
  )
  return res.data.sessionID
}

export async function getEvents(): Promise<MegaEvent[]> {
  const sessionId = localStorage?.getItem('sessionId')
  if (!sessionId) {
    return []
  }
  const res = await axios.get<EventResponse>(
    BACKEND_HOST + '/api/mega_events',
    {
      headers: {
        'X-Session': sessionId,
      },
    }
  )

  return res.data?.events || []
}

export async function getInscribedEvents(): Promise<number[]> {
  const sessionId = localStorage?.getItem('sessionId')
  if (!sessionId) {
    return []
  }
  const res = await axios.get<{ inscribedClasses: number[] }>(
    BACKEND_HOST + '/api/mega_inscribe',
    {
      headers: {
        'X-Session': sessionId,
      },
    }
  )

  return res.data?.inscribedClasses || []
}

export async function inscribeToEvent(classId: number): Promise<boolean> {
  const sessionId = localStorage?.getItem('sessionId')
  if (!sessionId) {
    return false
  }
  const res = await axios.post<{ success: boolean }>(
    BACKEND_HOST + '/api/mega_inscribe/' + classId,
    {},
    {
      headers: {
        'X-Session': sessionId,
      },
    }
  )

  return !!res.data?.success
}

export async function getMegaEventsBooked(
  classIds: number[]
): Promise<number[]> {
  const sessionId = localStorage?.getItem('sessionId')
  if (!sessionId) {
    return []
  }
  const res = await axios.get<number[]>(
    BACKEND_HOST + '/api/mega_events_booked?classIds=' + classIds.join(','),
    {
      headers: {
        'X-Session': sessionId,
      },
    }
  )

  return res.data || []
}
