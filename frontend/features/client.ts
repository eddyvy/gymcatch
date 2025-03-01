import axios from 'axios'

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
