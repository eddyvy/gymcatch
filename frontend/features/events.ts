export interface EventResponse {
  events: MegaEvent[]
}

export type MegaEvent = {
  activity_id: number
  activity_name: string
  affected_by_replacements: boolean
  attendees: number
  booking_info: BookingInfo
  booking_waiting_list: boolean
  booking_waiting_list_info: BookingWaitingListInfo
  bookings_appendable: boolean
  bookings_cancelable: boolean
  bookings_editable: boolean
  bookings_listable: boolean
  categories_ids: number[]
  color: string
  conflict: any[]
  duration: number
  end: string
  hour: string
  id: string
  instructors?: Instructor[]
  mobile: Mobile
  replacer: any
  resourceId: number
  room: string
  room_obj: RoomObj
  rotating: boolean
  session_id: number
  special_class?: string
  start: string
  startEditable: boolean
  substitution_instructors?: SubstitutionInstructor[]
  target: number
  title: string
  type_of_class?: string
  wday: number
}

export interface BookingInfo {
  available: boolean
  i_have_booked: boolean
  login_required: boolean
  pass_required: boolean
  places: Places
  products_info: ProductsInfo
  sold_out: boolean
  too_late: boolean
  too_soon: boolean
}

export interface Places {
  booked: number
  total: number
}

export interface ProductsInfo {}

export interface BookingWaitingListInfo {}

export interface Instructor {
  activities_ids: number[]
  avatar: string
  id: number
  name: string
  surname: string
}

export interface Mobile {
  color: string
  duration: string
  icon: string
  month_day: string
  roomOrder: any
  start_time: string
  subtitle: string
  title: string
  week_day: string
  week_day_short: string
}

export interface RoomObj {
  icon: Icon
  id: number
}

export interface Icon {
  css_class: string
}

export interface SubstitutionInstructor {
  avatar: string
  id: number
  name: string
  surname: string
}
