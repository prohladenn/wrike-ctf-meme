export interface TemplateDto {
  id: number
  name: string
  created_at: Date
  owner_id: number
  private_info?: string
}

export interface CreateTemplateDto {
  name: string
  comment: string
}
