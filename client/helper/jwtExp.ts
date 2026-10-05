/**
 * Decodes a JWT payload without verifying its signature.
 *
 * Safe to call only on a token that has already been proven genuine — the
 * server verifies the signature on every authenticated request, so this is
 * just a way to read claims out of a token we were handed. Never use the
 * result to make an authorization decision.
 */
export interface JwtClaims {
  user_id?: number;
  email?: string;
  role?: string;
  exp?: number;
  iat?: number;
}

export function decodeJwtClaims(token: string): JwtClaims | null {
  try {
    const [, payloadBase64] = token.split(".");
    if (!payloadBase64) return null;

    const decoded = JSON.parse(
      Buffer.from(payloadBase64, "base64").toString("utf-8"),
    );

    return decoded as JwtClaims;
  } catch {
    return null;
  }
}

export function getTokenRemainingSeconds(token: string): number {
  try {
    const decodedPayload = decodeJwtClaims(token);
    if (!decodedPayload?.exp) return 60 * 60 * 24; // Default fallback

    const currentTimeInSeconds = Math.floor(Date.now() / 1000);
    const remainingSeconds = decodedPayload.exp - currentTimeInSeconds;

    return Math.max(remainingSeconds, 0);
  } catch {
    return 60 * 60 * 24;
  }
}
