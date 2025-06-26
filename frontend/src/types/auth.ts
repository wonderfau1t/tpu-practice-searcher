export interface AuthResponse {
  result: AuthResult;
}

export interface AuthResult {
  role: 'student' | 'moderator' | 'headHR' | 'HR' | 'guest';
  accessToken: string;
  status: string;
}