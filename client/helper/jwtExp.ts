export function getTokenRemainingSeconds(token: string): number {
  try {
    const payloadBase64 = token.split(".")[1];
    const decodedPayload = JSON.parse(
      Buffer.from(payloadBase64, "base64").toString("utf-8"),
    );

    if (!decodedPayload.exp) return 60 * 60 * 24; // Default fallback

    const currentTimeInSeconds = Math.floor(Date.now() / 1000);
    const remainingSeconds = decodedPayload.exp - currentTimeInSeconds;

    return Math.max(remainingSeconds, 0);
  } catch {
    return 60 * 60 * 24;
  }
}
