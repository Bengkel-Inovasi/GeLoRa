import type { IAuthRepository, SignInPayload, AuthTokens } from '../../../ports/outbound/repository/auth';
import { httpClient } from '../client/http';

class AuthRepository implements IAuthRepository {
  signIn(payload: SignInPayload): Promise<AuthTokens> {
    return httpClient.post<AuthTokens>('/auth/sign-in', payload);
  }

  signUp(payload: { name: string; username: string; password: string }): Promise<AuthTokens> {
    return httpClient.post<AuthTokens>('/auth/sign-up', payload);
  }

  refresh(refreshToken: string): Promise<AuthTokens> {
    return httpClient.post<AuthTokens>('/auth/refresh', { refresh_token: refreshToken });
  }
}

export const authRepository = new AuthRepository();
