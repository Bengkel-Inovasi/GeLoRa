export interface SignInPayload {
  username: string;
  password: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

export interface IAuthRepository {
  signIn(payload: SignInPayload): Promise<AuthTokens>;
  signUp(payload: { name: string; username: string; password: string }): Promise<AuthTokens>;
  refresh(refreshToken: string): Promise<AuthTokens>;
}
