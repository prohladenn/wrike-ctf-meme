export interface MemeDto {
  id: number
  name: string
  created_at: Date
  owner_id: number
  owner_username: string
  template_id: number
}

export interface CreateMemeDto {
  template_id: string
  name: string
  caption: string
  color: string
}
