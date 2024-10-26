export interface LoginDto {
  username: string
  password: string
}

// tslint:disable-next-line:no-empty-interface
export interface RegisterDto extends LoginDto {}
